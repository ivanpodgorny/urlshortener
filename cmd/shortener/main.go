package main

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/ivanpodgorny/urlshortener/cmd/shortener/handler"
	"github.com/ivanpodgorny/urlshortener/cmd/shortener/storage"
	"log"
	"net/http"
)

func main() {
	var (
		r = chi.NewRouter()
		s = NewShortener(storage.NewMemory())
		h = handler.NewShortenURL(s)
	)

	r.Use(middleware.Recoverer)

	r.Post("/", h.Create)
	r.Get("/{id:[A-Za-z0-9_-]+}", h.Get)
	r.Post("/api/shorten", h.CreateJSON)

	log.Fatal(http.ListenAndServe(":8080", r))
}
