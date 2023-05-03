package storage

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"

	inerr "github.com/ivanpodgorny/urlshortener/internal/app/errors"
)

// Pg реализует интерфейс service.Storage для хранения url в PostgreSQL.
type Pg struct {
	db *sql.DB
}

// NewPg возвращает указатель на новый экземпляр Pg.
func NewPg(db *sql.DB) *Pg {
	return &Pg{db: db}
}

// Add сохраняет URL. Если URL был сохранен ранее, возвращает его id.
func (p *Pg) Add(ctx context.Context, id string, url string, userID string) (string, error) {
	_, err := p.db.ExecContext(ctx, "insert into urls (user_id, url_id, url) values ($1, $2, $3)", userID, id, url)

	if err != nil && err.(*pgconn.PgError).Code == pgerrcode.UniqueViolation {
		storedID := ""
		err = p.db.QueryRowContext(ctx, "select url_id from urls where url = $1", url).Scan(&storedID)

		return storedID, err
	}

	return id, err
}

// Get возвращает сохраненный URL по id. Если URL был помечен удаленным, возвращает
// ошибку errors.ErrURLIsDeleted.
func (p *Pg) Get(ctx context.Context, id string) (string, error) {
	var (
		url     = ""
		deleted = false
	)
	if err := p.db.QueryRowContext(ctx, "select url, deleted from urls where url_id = $1", id).Scan(&url, &deleted); err != nil {
		return url, err
	}

	if deleted {
		return url, inerr.ErrURLIsDeleted
	}

	return url, nil
}

// GetAllUser возвращает все сохраненные URL пользователя.
func (p *Pg) GetAllUser(ctx context.Context, userID string) map[string]string {
	data := map[string]string{}
	rows, err := p.db.QueryContext(ctx, "select url_id, url from urls where user_id = $1 and deleted = false", userID)
	if err != nil {
		return data
	}

	defer func(rows *sql.Rows) {
		_ = rows.Close()
	}(rows)

	for rows.Next() {
		var (
			id  = ""
			url = ""
		)
		err = rows.Scan(&id, &url)
		if err != nil {
			continue
		}

		data[id] = url
	}

	if err := rows.Err(); err != nil {
		return map[string]string{}
	}

	return data
}

// DeleteBatch удаляет URL с заданными id.
func (p *Pg) DeleteBatch(ctx context.Context, urlIDs []string, userID string) error {
	var (
		params       = make([]any, len(urlIDs)+1)
		placeholders = strings.Builder{}
	)
	params[0] = userID
	for i, urlID := range urlIDs {
		if i != 0 {
			placeholders.WriteString(",")
		}
		placeholders.WriteString(fmt.Sprintf("$%d", i+2))
		params[i+1] = urlID
	}

	_, err := p.db.ExecContext(ctx, `
update urls
set deleted = true
from (select unnest(array[`+placeholders.String()+`]) as url_id) as id_table
where user_id = $1
  and urls.url_id = id_table.url_id
	`, params...)

	return err
}
