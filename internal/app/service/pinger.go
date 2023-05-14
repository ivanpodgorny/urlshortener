package service

import (
	"context"
	"database/sql"
	"time"
)

// Pinger реализует метод для проверки доступности БД.
type Pinger struct {
	db *sql.DB
}

// NewPinger возвращает указатель на новый экземпляр Pinger.
func NewPinger(db *sql.DB) *Pinger {
	return &Pinger{db: db}
}

// Ping проверяет доступность БД.
func (p Pinger) Ping(ctx context.Context) bool {
	ctx, cancel := context.WithTimeout(ctx, 1*time.Second)
	defer cancel()

	err := p.db.PingContext(ctx)

	return err == nil
}
