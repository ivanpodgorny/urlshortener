package handler

import (
	"context"
	"github.com/go-chi/chi/v5"
	"io"
	"net/http"
	"net/url"
)

type ShortenURL struct {
	authenticator IdentityProvider
	shortener     Shortener
	baseURL       string
}

type IdentityProvider interface {
	UserIdentifier(r *http.Request) (string, error)
}

type Shortener interface {
	Shorten(ctx context.Context, url string, userID string) (string, error)
	Get(ctx context.Context, id string) (string, error)
	GetAllUser(ctx context.Context, userID string) map[string]string
}

func NewShortenURL(a IdentityProvider, s Shortener, b string) *ShortenURL {
	return &ShortenURL{
		authenticator: a,
		shortener:     s,
		baseURL:       b,
	}
}

// Create обрабатывает запрос на создание сокращенного URL.
// Оригинальный URL передается в теле запроса. В теле ответа приходит сокращенный URL.
func (h ShortenURL) Create(w http.ResponseWriter, r *http.Request) {
	userID, err := h.authenticator.UserIdentifier(r)
	if err != nil {
		unauthorized(w)

		return
	}

	b, err := io.ReadAll(r.Body)
	if err != nil || !h.isURL(string(b)) {
		badRequest(w)

		return
	}

	id, err := h.shortener.Shorten(r.Context(), string(b), userID)
	if err != nil {
		serverError(w)

		return
	}

	responseAsText(w, []byte(h.prepareShortenURL(id)), http.StatusCreated)
}

// CreateJSON обрабатывает запрос на создание сокращенного URL.
// Оригинальный URL передается в теле запроса в формате JSON {"url":"<some_url>"}.
// В теле ответа приходит JSON формата {"result":"<shorten_url>"} с сокращенным URL.
func (h ShortenURL) CreateJSON(w http.ResponseWriter, r *http.Request) {
	userID, err := h.authenticator.UserIdentifier(r)
	if err != nil {
		unauthorized(w)

		return
	}

	req := struct {
		URL string `json:"url"`
	}{}
	err = readJSONBody(&req, r)
	if err != nil || !h.isURL(req.URL) {
		badRequest(w)

		return
	}

	id, err := h.shortener.Shorten(r.Context(), req.URL, userID)
	if err != nil {
		serverError(w)

		return
	}

	responseAsJSON(
		w,
		struct {
			URL string `json:"result"`
		}{
			URL: h.prepareShortenURL(id),
		},
		http.StatusCreated,
	)
}

// Get обрабатывает запрос на получение оригинального URL из сокращенного.
// Возвращает ответ с кодом 307 и оригинальным URL в HTTP-заголовке Location.
func (h ShortenURL) Get(w http.ResponseWriter, r *http.Request) {
	u, err := h.shortener.Get(r.Context(), chi.URLParam(r, "id"))
	if err != nil {
		http.NotFound(w, r)

		return
	}

	redirect(w, u, http.StatusTemporaryRedirect)
}

// GetAllByCurrentUser возвращает все сокращенные URL пользователя, выполнившего запрос, в формате
// [{"short_url": "http://...", "original_url": "http://..."}, ...].
func (h ShortenURL) GetAllByCurrentUser(w http.ResponseWriter, r *http.Request) {
	userID, err := h.authenticator.UserIdentifier(r)
	if err != nil {
		unauthorized(w)

		return
	}

	type urlData struct {
		ShortURL    string `json:"short_url"`
		OriginalURL string `json:"original_url"`
	}
	resp := make([]urlData, 0)
	urls := h.shortener.GetAllUser(r.Context(), userID)
	for id, u := range urls {
		resp = append(resp, urlData{
			ShortURL:    h.prepareShortenURL(id),
			OriginalURL: u,
		})
	}

	if len(resp) == 0 {
		w.WriteHeader(http.StatusNoContent)

		return
	}

	responseAsJSON(w, resp, http.StatusOK)
}

func (h ShortenURL) isURL(str string) bool {
	u, err := url.Parse(str)

	return err == nil && u.Scheme != "" && u.Host != ""
}

func (h ShortenURL) prepareShortenURL(id string) string {
	return h.baseURL + "/" + id
}
