package service

import (
	"context"
	"database/sql"
	"time"
)

type Pinger struct {
	db *sql.DB
}

func NewPinger(db *sql.DB) *Pinger {
	return &Pinger{db: db}
}

func (p Pinger) Ping(ctx context.Context) bool {
	ctx, cancel := context.WithTimeout(ctx, 1*time.Second)
	defer cancel()

	err := p.db.PingContext(ctx)

	return err == nil
}
