package handler

import (
	"context"
	"net/http"
)

// Database реализует хендлеры для работы с БД.
type Database struct {
	pinger Pinger
}

// Pinger интерфейс сервиса проверки доступности БД.
type Pinger interface {
	Ping(ctx context.Context) bool
}

// NewDatabase возвращает указатель на новый экземпляр Database.
func NewDatabase(p Pinger) *Database {
	return &Database{
		pinger: p,
	}
}

// Ping обрабатывает запрос на проверку соединения с БД.
func (h Database) Ping(w http.ResponseWriter, r *http.Request) {
	if !h.pinger.Ping(r.Context()) {
		serverError(w)

		return
	}

	w.WriteHeader(http.StatusOK)
}
