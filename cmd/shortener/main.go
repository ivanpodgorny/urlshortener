package main

import (
	"compress/flate"
	"context"
	"crypto/tls"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/ivanpodgorny/urlshortener/internal/app/interceptor"

	"google.golang.org/grpc"

	"github.com/ivanpodgorny/urlshortener/internal/proto"

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
	if err := Execute(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatal(err)
	}
}

// Execute запускает веб-сервер.
// Для конфигурирования используются флаги и переменные окружения. Приоритет отдается
// значениям, заданным в переменных окружения.
// В качестве хранилища данных используется PostgreSQL, если указано DSN, иначе данные
// хранятся в файле на диске.
func Execute() error {
	cfg, err := config.NewBuilder().
		LoadFile().
		LoadFlags().
		LoadEnv().
		Build()
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
		r  = chi.NewRouter()
		cp = security.NewHMACTokenCreatorParser(cfg.HMACKey())
		a  = security.NewAuthenticator(
			security.NewCookieTokenStorage(cp),
			&security.RequestContextUserProvider{},
		)
		ga = security.NewGRPCAuthenticator(cp, security.NewGRPCContextUserProvider())
		wg = &sync.WaitGroup{}
		ss = service.NewShortener(store)
		sh = handler.NewShortenURL(a, ss, cfg.BaseURL(), wg)
		dh = handler.NewDatabase(service.NewPinger(db))
	)

	go func() {
		if err = startGRPCServer(cfg, ss, ga); err != nil {
			log.Printf("GRPC server error: %v", err)
		}
	}()

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
	r.With(middleware.Internal(cfg.TrustedSubnet())).Get("/api/internal/stats", sh.GetStat)
	r.Get("/ping", dh.Ping)

	fmt.Printf(buildInfo, buildVersion, buildDate, buildCommit)

	shutdownDone, err := startHTTPServer(cfg, r)
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	<-shutdownDone
	wg.Wait()

	log.Println("Server gracefully shutdown")

	return err
}

func startHTTPServer(cfg *config.Config, r *chi.Mux) (<-chan struct{}, error) {
	var (
		srv = &http.Server{
			Addr:    cfg.ServerAddress(),
			Handler: r,
		}
		sigCh          = make(chan os.Signal, 1)
		shutdown       = make(chan struct{})
		err      error = nil
	)

	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	go func() {
		<-sigCh
		log.Println("Starting server graceful shutdown...")
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
		defer cancel()

		if serr := srv.Shutdown(ctx); serr != nil {
			log.Printf("Error while server gracefully shutdown: %v", serr)
			if serr := srv.Close(); serr != nil {
				log.Printf("Error while server force shutdown: %v", serr)
			}
		}
		close(shutdown)
	}()

	if !cfg.EnableHTTPS() {
		err = srv.ListenAndServe()
	} else {
		cert, cerr := security.CreateCertificate()
		if cerr != nil {
			return shutdown, cerr
		}
		l, cerr := tls.Listen(
			"tcp",
			cfg.ServerAddress(),
			&tls.Config{
				Certificates: []tls.Certificate{cert},
			},
		)
		if cerr != nil {
			return shutdown, cerr
		}
		defer func(l net.Listener) {
			err = l.Close()
		}(l)

		err = srv.Serve(l)
	}

	return shutdown, err
}

func startGRPCServer(cfg *config.Config, s handler.Shortener, a *security.GRPCAuthenticator) error {
	listen, err := net.Listen("tcp", cfg.GRPCServerAddress())
	if err != nil {
		return err
	}

	gs := grpc.NewServer(grpc.UnaryInterceptor(interceptor.Authenticate(a)))
	proto.RegisterShortenerServer(gs, handler.NewShortenerGRPCServer(a, s))

	return gs.Serve(listen)
}
