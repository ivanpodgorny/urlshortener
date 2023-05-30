package main

import (
	"compress/flate"
	"crypto/tls"
	"database/sql"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
	_ "github.com/jackc/pgx/v5/stdlib"

	"github.com/ivanpodgorny/urlshortener/internal/app/config"
	"github.com/ivanpodgorny/urlshortener/internal/app/handler"
	"github.com/ivanpodgorny/urlshortener/internal/app/middleware"
	"github.com/ivanpodgorny/urlshortener/internal/app/migrations"
	"github.com/ivanpodgorny/urlshortener/internal/app/security"
	"github.com/ivanpodgorny/urlshortener/internal/app/service"
	"github.com/ivanpodgorny/urlshortener/internal/app/storage"
)

const buildInfo = "Build version: %s\nBuild date: %s\nBuild commit: %s\n"

var (
	buildVersion = "N/A"
	buildDate    = "N/A"
	buildCommit  = "N/A"
)

func main() {
	log.Fatal(Execute())
}

// Execute запускает веб-сервер.
// Для конфигурирования используются флаги и переменные окружения. Приоритет отдается
// значениям, заданным в переменных окружения.
// В качестве хранилища данных используется PostgreSQL, если указано DSN, иначе данные
// хранятся в файле на диске.
func Execute() error {
	cfg, err := config.NewBuilder().LoadFlags().LoadEnv().Build()
	if err != nil {
		return err
	}

	var file *os.File
	if cfg.FileStoragePath() != "" {
		file, err = os.OpenFile(cfg.FileStoragePath(), os.O_RDWR|os.O_CREATE, 0600)
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
		if err = migrations.Up(db); err != nil {
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

	r.Mount("/debug", chimiddleware.Profiler())
	r.Post("/", sh.Create)
	r.Get("/{id:[A-Za-z0-9_-]+}", sh.Get)
	r.Post("/api/shorten", sh.CreateJSON)
	r.Post("/api/shorten/batch", sh.CreateBatch)
	r.Get("/api/user/urls", sh.GetAllByCurrentUser)
	r.Delete("/api/user/urls", sh.DeleteBatch)
	r.Get("/ping", dh.Ping)

	fmt.Printf(buildInfo, buildVersion, buildDate, buildCommit)

	if !cfg.EnableHTTPS() {
		return http.ListenAndServe(cfg.ServerAddress(), r)
	}

	srv := &http.Server{
		Handler: r,
	}

	cert, err := security.CreateCertificate()
	if err != nil {
		return err
	}
	l, err := tls.Listen(
		"tcp",
		cfg.ServerAddress(),
		&tls.Config{
			Certificates: []tls.Certificate{cert},
		},
	)
	if err != nil {
		return err
	}
	defer func(l net.Listener) {
		err = l.Close()
	}(l)

	err = srv.Serve(l)

	return err
}
