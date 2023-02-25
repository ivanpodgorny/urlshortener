package handler

import (
	"context"
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"io"
	"net/http"
	"net/url"
)

type ShortenURL struct {
	shortener Shortener
	baseURL   string
}

type Shortener interface {
	Shorten(ctx context.Context, url string) (string, error)
	Get(ctx context.Context, id string) (string, error)
}

func NewShortenURL(s Shortener, b string) *ShortenURL {
	return &ShortenURL{
		shortener: s,
		baseURL:   b,
	}
}

func BadRequest(w http.ResponseWriter) {
	http.Error(w, "400 bad request", http.StatusBadRequest)
}

func ServerError(w http.ResponseWriter) {
	http.Error(w, "500 internal server error", http.StatusInternalServerError)
}

// Create обрабатывает запрос на создание сокращенного URL.
// Оригинальный URL передается в теле запроса. В теле ответа приходит сокращенный URL.
func (h ShortenURL) Create(w http.ResponseWriter, r *http.Request) {
	b, err := io.ReadAll(r.Body)
	if err != nil || !h.isURL(string(b)) {
		BadRequest(w)

		return
	}

	id, err := h.shortener.Shorten(r.Context(), string(b))
	if err != nil {
		ServerError(w)

		return
	}

	w.WriteHeader(http.StatusCreated)
	if _, err := w.Write([]byte(h.prepareShortenURL(id))); err != nil {
		ServerError(w)
	}
}

// CreateJSON обрабатывает запрос на создание сокращенного URL.
// Оригинальный URL передается в теле запроса в формате JSON {"url":"<some_url>"}.
// В теле ответа приходит JSON формата {"result":"<shorten_url>"} с сокращенным URL.
func (h ShortenURL) CreateJSON(w http.ResponseWriter, r *http.Request) {
	b, err := io.ReadAll(r.Body)
	if err != nil {
		BadRequest(w)

		return
	}

	req := struct {
		URL string `json:"url"`
	}{}
	err = json.Unmarshal(b, &req)
	if err != nil || !h.isURL(req.URL) {
		BadRequest(w)

		return
	}

	id, err := h.shortener.Shorten(r.Context(), req.URL)
	if err != nil {
		ServerError(w)

		return
	}

	resp := struct {
		URL string `json:"result"`
	}{
		URL: h.prepareShortenURL(id),
	}
	respJSON, err := json.Marshal(resp)
	if err != nil {
		ServerError(w)

		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if _, err := w.Write(respJSON); err != nil {
		ServerError(w)
	}
}

// Get обрабатывает запрос на получение оригинального URL из сокращенного.
// Возвращает ответ с кодом 307 и оригинальным URL в HTTP-заголовке Location.
func (h ShortenURL) Get(w http.ResponseWriter, r *http.Request) {
	u, err := h.shortener.Get(r.Context(), chi.URLParam(r, "id"))
	if err != nil {
		http.NotFound(w, r)

		return
	}

	w.Header().Set("Location", u)
	w.WriteHeader(http.StatusTemporaryRedirect)
}

func (h ShortenURL) isURL(str string) bool {
	u, err := url.Parse(str)

	return err == nil && u.Scheme != "" && u.Host != ""
}

func (h ShortenURL) prepareShortenURL(id string) string {
	return h.baseURL + "/" + id
}
