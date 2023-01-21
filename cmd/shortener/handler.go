package main

import (
	"fmt"
	"github.com/ivanpodgorny/urlshortener/cmd/shortener/router"
	"io"
	"net/http"
	"net/url"
)

type Handler struct {
	shortener *Shortener
}

func NewHandler(s *Shortener) *Handler {
	return &Handler{shortener: s}
}

func (h Handler) Create(w http.ResponseWriter, r *http.Request, _ ...string) {
	b, err := io.ReadAll(r.Body)
	if err == nil && h.isUrl(string(b[:])) {
		id, err := h.shortener.Shorten(r.Context(), string(b[:]))
		if err == nil {
			w.WriteHeader(http.StatusCreated)
			host := r.Host
			fmt.Println(host)
			_, err := w.Write([]byte("http://" + r.Host + "/" + id))
			if err != nil {
				router.BadRequest(w)
			}
		} else {
			router.BadRequest(w)
		}
	} else {
		router.BadRequest(w)
	}
}

func (h Handler) Get(w http.ResponseWriter, r *http.Request, params ...string) {
	u, err := h.shortener.Get(r.Context(), params[0])
	if err == nil {
		w.Header().Set("Location", u)
		w.WriteHeader(http.StatusTemporaryRedirect)
	} else {
		http.NotFound(w, r)
	}
}

func (h Handler) isUrl(str string) bool {
	u, err := url.Parse(str)

	return err == nil && u.Scheme != "" && u.Host != ""
}
