package storage

import (
	"context"
	"database/sql"
)

type Pg struct {
	db *sql.DB
}

func NewPg(db *sql.DB) *Pg {
	return &Pg{db: db}
}

func (p *Pg) Add(ctx context.Context, id string, url string, userID string) error {
	_, err := p.db.ExecContext(ctx, "insert into urls (user_id, url_id, url) values ($1, $2, $3)", userID, id, url)

	return err
}

func (p *Pg) Get(ctx context.Context, id string) (string, error) {
	url := ""
	err := p.db.QueryRowContext(ctx, "select url from urls where url_id = $1", id).Scan(&url)

	return url, err
}

func (p *Pg) GetAllUser(ctx context.Context, userID string) map[string]string {
	data := map[string]string{}
	rows, err := p.db.QueryContext(ctx, "select url_id, url from urls where user_id = $1", userID)
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
