package handler

import (
	"bytes"
	"context"
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

type ShortenerMock struct {
	mock.Mock
}

func (m *ShortenerMock) Shorten(_ context.Context, _ string) (string, error) {
	args := m.Called()

	return args.String(0), args.Error(1)
}

func (m *ShortenerMock) Get(_ context.Context, _ string) (string, error) {
	args := m.Called()

	return args.String(0), args.Error(1)
}

func TestShortenURLHandler_Create(t *testing.T) {
	var (
		urlID     = "1i-CBrzwyMkL"
		shortener = &ShortenerMock{}
	)

	shortener.
		On("Shorten").Return("", errors.New("")).Once().
		On("Shorten").Return(urlID, nil).Once()
	handler := ShortenURL{
		shortener: shortener,
	}

	tests := []struct {
		name             string
		body             io.Reader
		wantStatusCode   int
		wantResponseBody string
	}{
		{
			name:           "пустое тело запроса",
			wantStatusCode: http.StatusBadRequest,
		},
		{
			name:           "некорректный URL",
			body:           bytes.NewBuffer([]byte("file://")),
			wantStatusCode: http.StatusBadRequest,
		},
		{
			name:           "ошибка создания сокращенного URL",
			body:           bytes.NewBuffer([]byte("https://ya.ru/")),
			wantStatusCode: http.StatusInternalServerError,
		},
		{
			name:             "успешное выполнение запроса",
			body:             bytes.NewBuffer([]byte("https://ya.ru/")),
			wantStatusCode:   http.StatusCreated,
			wantResponseBody: "http://example.com/" + urlID,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := sendTestRequest(http.MethodPost, "/", tt.body, handler.Create)
			assert.Equal(t, tt.wantStatusCode, result.StatusCode)
			if tt.wantResponseBody != "" {
				b, err := io.ReadAll(result.Body)
				require.NoError(t, err)
				assert.Equal(t, tt.wantResponseBody, string(b[:]))
			}
			err := result.Body.Close()
			require.NoError(t, err)
		})
	}
}

func TestShortenURLHandler_Get(t *testing.T) {
	var (
		url       = "https://ya.ru/"
		urlID     = "1i-CBrzwyMkL"
		shortener = &ShortenerMock{}
	)

	shortener.
		On("Get").Return("", errors.New("")).Once().
		On("Get").Return(url, nil).Once()
	handler := ShortenURL{
		shortener: shortener,
	}

	tests := []struct {
		name           string
		wantStatusCode int
		wantLocation   string
	}{
		{
			name:           "не найден URL для переданного ID",
			wantStatusCode: http.StatusNotFound,
		},
		{
			name:           "успешное выполнение запроса",
			wantStatusCode: http.StatusTemporaryRedirect,
			wantLocation:   url,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := sendTestRequest(http.MethodGet, "/"+urlID, nil, handler.Get)
			assert.Equal(t, tt.wantStatusCode, result.StatusCode)
			if tt.wantLocation != "" {
				assert.Equal(t, tt.wantLocation, result.Header.Get("Location"))
			}
			err := result.Body.Close()
			require.NoError(t, err)
		})
	}
}

func sendTestRequest(method string, target string, body io.Reader, handler http.HandlerFunc) *http.Response {
	request := httptest.NewRequest(method, target, body)
	w := httptest.NewRecorder()
	handler(w, request)

	return w.Result()
}
