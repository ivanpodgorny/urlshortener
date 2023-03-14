package main

import (
	"compress/flate"
	"database/sql"
	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/ivanpodgorny/urlshortener/internal/app/config"
	"github.com/ivanpodgorny/urlshortener/internal/app/handler"
	"github.com/ivanpodgorny/urlshortener/internal/app/middleware"
	"github.com/ivanpodgorny/urlshortener/internal/app/migrations"
	"github.com/ivanpodgorny/urlshortener/internal/app/security"
	"github.com/ivanpodgorny/urlshortener/internal/app/service"
	"github.com/ivanpodgorny/urlshortener/internal/app/storage"
	_ "github.com/jackc/pgx/v5/stdlib"
	"log"
	"net/http"
	"os"
)

func main() {
	log.Fatal(Execute())
}

func Execute() error {
	cfg, err := config.NewBuilder().LoadFlags().LoadEnv().Build()
	if err != nil {
		return err
	}

	var file *os.File
	if cfg.FileStoragePath() != "" {
		file, err = os.OpenFile(cfg.FileStoragePath(), os.O_RDWR|os.O_CREATE|os.O_APPEND, 0600)
		if err != nil {
			return err
		}

		defer func(file *os.File) {
			err = file.Close()
		}(file)
	}

	db, err := sql.Open("pgx", cfg.DatabaseDSN())
	if err != nil {
		return err
	}

	defer func(db *sql.DB) {
		err = db.Close()
	}(db)

	var store service.Storage
	store = storage.NewMemory(file)
	if cfg.DatabaseDSN() != "" {
		if err := migrations.Up(db); err != nil {
			return err
		}

		store = storage.NewPg(db)
	}

	var (
		r = chi.NewRouter()
		a = security.NewAuthenticator(
			security.NewCookieTokenStorage(security.NewHMACTokenCreatorParser(cfg.HMACKey())),
			&security.RequestContextUserProvider{},
		)
		sh = handler.NewShortenURL(a, service.NewShortener(store), cfg.BaseURL())
		dh = handler.NewDatabase(service.NewPinger(db))
	)

	r.Use(chimiddleware.Recoverer)
	r.Use(chimiddleware.Compress(flate.BestSpeed))
	r.Use(middleware.Decompress())
	r.Use(middleware.Authenticate(a))

	r.Post("/", sh.Create)
	r.Get("/{id:[A-Za-z0-9_-]+}", sh.Get)
	r.Post("/api/shorten", sh.CreateJSON)
	r.Post("/api/shorten/batch", sh.CreateBatch)
	r.Get("/api/user/urls", sh.GetAllByCurrentUser)
	r.Get("/ping", dh.Ping)

	err = http.ListenAndServe(cfg.ServerAddress(), r)

	return err
}
