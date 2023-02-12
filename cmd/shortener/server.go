package main

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/ivanpodgorny/urlshortener/cmd/shortener/handler"
	"github.com/ivanpodgorny/urlshortener/cmd/shortener/service"
	"github.com/ivanpodgorny/urlshortener/cmd/shortener/storage"
	"net/http"
)

func Run() error {
	cfg, err := LoadConfig()
	if err != nil {
		return err
	}

	var (
		r = chi.NewRouter()
		s = service.NewShortener(storage.NewMemory())
		h = handler.NewShortenURL(s, cfg.BaseURL)
	)

	r.Use(middleware.Recoverer)

	r.Post("/", h.Create)
	r.Get("/{id:[A-Za-z0-9_-]+}", h.Get)
	r.Post("/api/shorten", h.CreateJSON)

	return http.ListenAndServe(cfg.ServerAddress, r)
}
