package handler

import (
	"context"
	"github.com/ivanpodgorny/urlshortener/cmd/shortener/router"
	"io"
	"net/http"
	"net/url"
)

type ShortenURL struct {
	shortener Shortener
}

type Shortener interface {
	Shorten(ctx context.Context, url string) (string, error)
	Get(ctx context.Context, id string) (string, error)
}

func NewShortenURL(s Shortener) *ShortenURL {
	return &ShortenURL{shortener: s}
}

// Create обрабатывает запрос на создание сокращенного URL.
// Оригинальный URL передается в теле запроса. В теле ответа приходит сокращенный URL.
func (h ShortenURL) Create(w http.ResponseWriter, r *http.Request, _ ...string) {
	b, err := io.ReadAll(r.Body)
	if err == nil && h.isURL(string(b[:])) {
		id, err := h.shortener.Shorten(r.Context(), string(b[:]))
		if err == nil {
			w.WriteHeader(http.StatusCreated)
			if _, err := w.Write([]byte("http://" + r.Host + "/" + id)); err != nil {
				router.ServerError(w)
			}
		} else {
			router.ServerError(w)
		}
	} else {
		router.BadRequest(w)
	}
}

// Get обрабатывает запрос на получение оригинального URL из сокращенного.
// Возвращает ответ с кодом 307 и оригинальным URL в HTTP-заголовке Location.
func (h ShortenURL) Get(w http.ResponseWriter, r *http.Request, params ...string) {
	u, err := h.shortener.Get(r.Context(), params[0])
	if err == nil {
		w.Header().Set("Location", u)
		w.WriteHeader(http.StatusTemporaryRedirect)
	} else {
		http.NotFound(w, r)
	}
}

func (h ShortenURL) isURL(str string) bool {
	u, err := url.Parse(str)

	return err == nil && u.Scheme != "" && u.Host != ""
}
