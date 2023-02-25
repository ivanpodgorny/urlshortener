package main

import (
	"compress/flate"
	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/ivanpodgorny/urlshortener/internal/app/handler"
	"github.com/ivanpodgorny/urlshortener/internal/app/middleware"
	"github.com/ivanpodgorny/urlshortener/internal/app/security"
	"github.com/ivanpodgorny/urlshortener/internal/app/service"
	"github.com/ivanpodgorny/urlshortener/internal/app/storage"
	"net/http"
)

func Run() error {
	cfg, err := LoadConfig()
	if err != nil {
		return err
	}

	var store service.Storage
	store = storage.NewMemory()
	if cfg.FileStoragePath != "" {
		store, err = storage.NewFile(cfg.FileStoragePath)
		if err != nil {
			return err
		}

		defer func(s *storage.File) {
			err = s.Close()
		}(store.(*storage.File))
	}

	var (
		r = chi.NewRouter()
		s = service.NewShortener(store)
		a = security.NewAuthenticator(
			security.NewCookieTokenStorage(security.NewHMACTokenCreatorParser(cfg.HMACKey)),
			security.RequestContextUserProvider{},
		)
		h = handler.NewShortenURL(a, s, cfg.BaseURL)
	)

	r.Use(chimiddleware.Recoverer)
	r.Use(chimiddleware.Compress(flate.BestSpeed))
	r.Use(middleware.Decompress())
	r.Use(middleware.Authenticate(a))

	r.Post("/", h.Create)
	r.Get("/{id:[A-Za-z0-9_-]+}", h.Get)
	r.Post("/api/shorten", h.CreateJSON)
	r.Get("/api/user/urls", h.GetAllByCurrentUser)

	err = http.ListenAndServe(cfg.ServerAddress, r)

	return err
}
