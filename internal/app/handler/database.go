package handler

import (
	"context"
	"net/http"
)

type Database struct {
	pinger Pinger
}

type Pinger interface {
	Ping(ctx context.Context) bool
}

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
