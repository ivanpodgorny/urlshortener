package handler

import (
	"context"
	"errors"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"

	inerr "github.com/ivanpodgorny/urlshortener/internal/app/errors"
	"github.com/ivanpodgorny/urlshortener/internal/app/validator"
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
	Shorten(ctx context.Context, url string, userID string) (string, bool, error)
	Get(ctx context.Context, id string) (string, error)
	GetAllUser(ctx context.Context, userID string) map[string]string
	DeleteBatch(ctx context.Context, urlIDs []string, userID string) error
}

const deleteBatchSize = 250

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
	if err != nil || !h.validateURL(string(b)) {
		badRequest(w)

		return
	}

	id, inserted, err := h.shortener.Shorten(r.Context(), string(b), userID)
	if err != nil {
		serverError(w)

		return
	}

	status := http.StatusCreated
	if !inserted {
		status = http.StatusConflict
	}

	responseAsText(w, []byte(h.prepareShortenURL(id)), status)
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
	if err != nil || !h.validateURL(req.URL) {
		badRequest(w)

		return
	}

	id, inserted, err := h.shortener.Shorten(r.Context(), req.URL, userID)
	if err != nil {
		serverError(w)

		return
	}

	status := http.StatusCreated
	if !inserted {
		status = http.StatusConflict
	}

	responseAsJSON(
		w,
		struct {
			URL string `json:"result"`
		}{
			URL: h.prepareShortenURL(id),
		},
		status,
	)
}

// CreateBatch обрабатывает запрос на создание нескольких сокращенных URL.
// Оригинальные URL передаются в теле запроса в формате JSON
// [{"correlation_id": "<строковый идентификатор>", "original_url": "<URL для сокращения>"}, ...].
// В теле ответа приходит JSON формата
// [{"correlation_id": "<строковый идентификатор>", "short_url": "<сокращённый URL>"}, ... ]  с сокращенными URL.
func (h ShortenURL) CreateBatch(w http.ResponseWriter, r *http.Request) {
	userID, err := h.authenticator.UserIdentifier(r)
	if err != nil {
		unauthorized(w)

		return
	}

	type origBatchItem struct {
		ID  string `json:"correlation_id"`
		URL string `json:"original_url"`
	}
	req := make([]origBatchItem, 0)
	if err = readJSONBody(&req, r); err != nil {
		badRequest(w)

		return
	}

	if valid, _ := validator.Validate[[]origBatchItem](req, validator.Size[origBatchItem](1000)); !valid {
		badRequest(w)

		return
	}

	type shortenBatchItem struct {
		ID  string `json:"correlation_id"`
		URL string `json:"short_url"`
	}
	resp := make([]shortenBatchItem, 0, len(req))
	for _, u := range req {
		if !h.validateURL(u.URL) {
			continue
		}

		id, _, err := h.shortener.Shorten(r.Context(), u.URL, userID)
		if err != nil {
			continue
		}

		resp = append(resp, shortenBatchItem{
			ID:  u.ID,
			URL: h.prepareShortenURL(id),
		})
	}

	responseAsJSON(w, resp, http.StatusCreated)
}

// Get обрабатывает запрос на получение оригинального URL из сокращенного.
// Возвращает ответ с кодом 307 и оригинальным URL в HTTP-заголовке Location.
// Если URL был удален пользователем, возвращает ответ с кодом 410.
func (h ShortenURL) Get(w http.ResponseWriter, r *http.Request) {
	u, err := h.shortener.Get(r.Context(), chi.URLParam(r, "id"))
	if errors.Is(err, inerr.ErrURLIsDeleted) {
		w.WriteHeader(http.StatusGone)

		return
	}

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
	urls := h.shortener.GetAllUser(r.Context(), userID)
	resp := make([]urlData, 0, len(urls))
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

// DeleteBatch принимает список идентификаторов сокращённых URL для удаления в формате
// [ "a", "b", "c", "d", ...]. В случае успешного приема запроса возвращает ответ с кодом 202.
func (h ShortenURL) DeleteBatch(w http.ResponseWriter, r *http.Request) {
	userID, err := h.authenticator.UserIdentifier(r)
	if err != nil {
		unauthorized(w)

		return
	}

	urlIDs := make([]string, 0)
	if err = readJSONBody(&urlIDs, r); err != nil {
		badRequest(w)

		return
	}

	idsCount := len(urlIDs)
	for i := 0; i < idsCount; i += deleteBatchSize {
		end := i + deleteBatchSize
		if end > idsCount {
			end = idsCount
		}

		go func(chunk []string) {
			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()

			if err = h.shortener.DeleteBatch(ctx, chunk, userID); err != nil {
				log.Println("ошибка удаления url: " + err.Error())
			}
		}(urlIDs[i:end])
	}

	w.WriteHeader(http.StatusAccepted)
}

func (h ShortenURL) validateURL(u string) bool {
	valid, _ := validator.Validate[string](u, validator.IsURL, validator.Length(2000))

	return valid
}

func (h ShortenURL) prepareShortenURL(id string) string {
	return h.baseURL + "/" + id
}
