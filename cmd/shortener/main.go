package main

import (
	"github.com/ivanpodgorny/urlshortener/cmd/shortener/router"
	"github.com/ivanpodgorny/urlshortener/cmd/shortener/storage"
	"log"
	"net/http"
)

func main() {
	var (
		r = router.NewRouter()
		s = NewShortener(storage.NewMemory())
		h = NewHandler(s)
	)

	r.Add(http.MethodPost, "/", h.Create)
	r.Add(http.MethodGet, `/[A-Za-z0-9_-]+`, h.Get)

	log.Fatal(http.ListenAndServe(":8080", r))
}
