package storage

import (
	"context"
	"database/sql"
	"errors"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/DATA-DOG/go-txdb"
	"github.com/ivanpodgorny/urlshortener/internal/app/config"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestPgRealConnection(t *testing.T) {
	var (
		id                = "id1"
		wrongID           = "id2"
		url               = "https://ya.ru/"
		userID            = "438c4b98-fc98-45cf-ac63-c4a86fbd4ff4"
		userWithoutURLsID = "02872d15-5047-406c-a989-ee1b07465169"
		ctx               = context.Background()
		db, err           = setupTestDB(t)
		s                 = NewPg(db)
	)

	if err != nil {
		return
	}

	defer func(db *sql.DB) {
		_ = db.Close()
	}(db)

	insertedID, err := s.Add(ctx, id, url, userID)
	assert.NoError(t, err, "добавление новой записи")
	assert.Equal(t, id, insertedID, "добавление новой записи")
	stored, err := s.Get(ctx, id)
	assert.NoError(t, err, "получение записи")
	assert.Equal(t, url, stored, "получение записи")
	_, err = s.Get(ctx, wrongID)
	assert.Error(t, err, "получение несуществующей записи записи")
	urls := s.GetAllUser(ctx, userID)
	assert.Equal(t, map[string]string{id: url}, urls, "получение URL пользователя")
	urls = s.GetAllUser(ctx, userWithoutURLsID)
	assert.Equal(t, map[string]string{}, urls, "получение URL пользователя, не добавлявшего URL")
	_, err = s.Add(ctx, id, url, userID)
	assert.Error(t, err, "добавление записи c существующим id")
}

func TestPgUniqueUrl(t *testing.T) {
	var (
		ctx           = context.Background()
		url           = "https://ya.ru/"
		urlIDInserted = "fE2ZNnnhOuYG7oMi"
		urlIDExisted  = "6Qq362Ml98Y15zeb"
		userID        = "438c4b98-fc98-45cf-ac63-c4a86fbd4ff4"
	)

	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	require.NoError(t, err)
	s := NewPg(db)

	mock.ExpectExec("insert into urls (user_id, url_id, url) values ($1, $2, $3)").
		WithArgs(userID, urlIDInserted, url).
		WillReturnError(&pgconn.PgError{Code: pgerrcode.UniqueViolation})
	mock.ExpectQuery("select url_id from urls where url = $1").
		WithArgs(url).
		WillReturnRows(sqlmock.NewRows([]string{"url"}).AddRow(urlIDExisted))
	id, err := s.Add(ctx, urlIDInserted, url, userID)
	assert.NoError(t, err)
	assert.Equal(t, urlIDExisted, id)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func setupTestDB(t *testing.T) (*sql.DB, error) {
	cfg, err := config.NewBuilder().LoadEnv().Build()
	require.NoError(t, err)
	if cfg.DatabaseDSN() == "" {
		return nil, errors.New("no dsn")
	}

	txdb.Register("txdb", "pgx", cfg.DatabaseDSN())
	db, err := sql.Open("txdb", "identifier")
	require.NoError(t, err)

	return db, nil
}
