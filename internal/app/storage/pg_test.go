package storage

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"errors"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/DATA-DOG/go-txdb"
	"github.com/ivanpodgorny/urlshortener/internal/app/config"
	inerr "github.com/ivanpodgorny/urlshortener/internal/app/errors"
	"github.com/ivanpodgorny/urlshortener/internal/app/security"
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
		idToDelete        = "id3"
		url               = "https://ya.ru/"
		urlToDelete       = "https://www.google.com/"
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
	_, _ = s.Add(ctx, idToDelete, urlToDelete, userID)
	err = s.DeleteBatch(ctx, []string{idToDelete}, userWithoutURLsID)
	assert.NoError(t, err, "попытка удаления чужой записи")
	notDeletedURL, err := s.Get(ctx, idToDelete)
	assert.NoError(t, err, "попытка удаления чужой записи")
	assert.Equal(t, urlToDelete, notDeletedURL, "попытка удаления чужой записи")
	err = s.DeleteBatch(ctx, []string{idToDelete}, userID)
	assert.NoError(t, err, "удаление записи")
	_, err = s.Get(ctx, idToDelete)
	assert.ErrorIs(t, err, inerr.ErrURLIsDeleted, "получение удаленной записи")
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

func BenchmarkPg_GetAllUser(b *testing.B) {
	var (
		db, mock, _ = sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
		s           = NewPg(db)
		ctx         = context.Background()
		userID      = "1"
		rows        = sqlmock.NewRows([]string{"url_id", "url"})
	)
	for i := 0; i < 1000; i++ {
		id, _ := security.GenerateRandomString(16)
		rows = rows.AddRow(id, "https://ya.ru/")
	}
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		mock.ExpectQuery("select url_id, url from urls where user_id = $1 and deleted = false").
			WithArgs(userID).
			WillReturnRows(rows)
		s.GetAllUser(ctx, userID)
	}
}

func BenchmarkPg_DeleteBatch(b *testing.B) {
	var (
		db, mock, _ = sqlmock.New()
		s           = NewPg(db)
		ctx         = context.Background()
		userID      = "1"
		urlIDs      = make([]string, 250)
		args        = make([]driver.Value, 251)
	)
	args[0] = userID
	for i := range urlIDs {
		id, _ := security.GenerateRandomString(16)
		urlIDs[i] = id
		args[i+1] = id
	}
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		mock.ExpectExec("update urls").
			WithArgs(args...).
			WillReturnResult(sqlmock.NewResult(0, 1))
		_ = s.DeleteBatch(ctx, urlIDs, userID)
	}
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
