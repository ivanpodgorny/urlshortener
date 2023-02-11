package handler

import (
	"bytes"
	"context"
	"encoding/json"
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

func TestShortenURLHandler_CreateSuccess(t *testing.T) {
	var (
		urlID     = "1i-CBrzwyMkL"
		shortener = &ShortenerMock{}
	)

	shortener.On("Shorten").Return(urlID, nil).Once()
	handler := ShortenURL{
		shortener: shortener,
	}

	result := sendTestRequest(http.MethodPost, "/", bytes.NewBuffer([]byte("https://ya.ru/")), handler.Create)
	assert.Equal(t, http.StatusCreated, result.StatusCode)
	b, err := io.ReadAll(result.Body)
	require.NoError(t, err)
	assert.Equal(t, "http://example.com/"+urlID, string(b[:]))
	err = result.Body.Close()
	require.NoError(t, err)
	shortener.AssertExpectations(t)
}

func TestShortenURLHandler_CreateWithErrors(t *testing.T) {
	shortener := &ShortenerMock{}
	shortener.On("Shorten").Return("", errors.New("")).Once()
	handler := ShortenURL{
		shortener: shortener,
	}

	tests := []struct {
		name           string
		body           io.Reader
		wantStatusCode int
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
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := sendTestRequest(http.MethodPost, "/", tt.body, handler.Create)
			assert.Equal(t, tt.wantStatusCode, result.StatusCode)
			err := result.Body.Close()
			require.NoError(t, err)
		})
	}
	shortener.AssertExpectations(t)
}

func TestShortenURLHandler_CreateJSONSuccess(t *testing.T) {
	var (
		urlID     = "1i-CBrzwyMkL"
		shortener = &ShortenerMock{}
	)

	shortener.On("Shorten").Return(urlID, nil).Once()
	handler := ShortenURL{
		shortener: shortener,
	}

	result := sendTestRequest(http.MethodPost, "/", bytes.NewBuffer([]byte(`{"url":"https://ya.ru/"}`)), handler.CreateJSON)
	assert.Equal(t, http.StatusCreated, result.StatusCode)
	assert.Equal(t, "application/json", result.Header.Get("Content-Type"))
	b, err := io.ReadAll(result.Body)
	require.NoError(t, err)
	resp := struct {
		URL string `json:"result"`
	}{}
	err = json.Unmarshal(b, &resp)
	require.NoError(t, err)
	assert.Equal(t, "http://example.com/"+urlID, resp.URL)
	err = result.Body.Close()
	require.NoError(t, err)
	shortener.AssertExpectations(t)
}

func TestShortenURLHandler_CreateJSONWithErrors(t *testing.T) {
	shortener := &ShortenerMock{}
	shortener.On("Shorten").Return("", errors.New("")).Once()
	handler := ShortenURL{
		shortener: shortener,
	}

	tests := []struct {
		name               string
		body               io.Reader
		wantStatusCode     int
		wantResponseResult string
	}{
		{
			name:           "неверный формат JSON",
			body:           bytes.NewBuffer([]byte(`{"u":"https://ya.ru/"}`)),
			wantStatusCode: http.StatusBadRequest,
		},
		{
			name:           "некорректный URL",
			body:           bytes.NewBuffer([]byte(`{"url":"file://"}`)),
			wantStatusCode: http.StatusBadRequest,
		},
		{
			name:           "ошибка создания сокращенного URL",
			body:           bytes.NewBuffer([]byte(`{"url":"https://ya.ru/"}`)),
			wantStatusCode: http.StatusInternalServerError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := sendTestRequest(http.MethodPost, "/", tt.body, handler.CreateJSON)
			assert.Equal(t, tt.wantStatusCode, result.StatusCode)
			err := result.Body.Close()
			require.NoError(t, err)
		})
	}
	shortener.AssertExpectations(t)
}

func TestShortenURLHandler_GetSuccess(t *testing.T) {
	var (
		url       = "https://ya.ru/"
		urlID     = "1i-CBrzwyMkL"
		shortener = &ShortenerMock{}
	)

	shortener.On("Get").Return(url, nil).Once()
	handler := ShortenURL{
		shortener: shortener,
	}

	result := sendTestRequest(http.MethodGet, "/"+urlID, nil, handler.Get)
	assert.Equal(t, http.StatusTemporaryRedirect, result.StatusCode)
	assert.Equal(t, url, result.Header.Get("Location"))
	err := result.Body.Close()
	require.NoError(t, err)
	shortener.AssertExpectations(t)
}

func TestShortenURLHandler_GetWithErrors(t *testing.T) {
	var (
		urlID     = "1i-CBrzwyMkL"
		shortener = &ShortenerMock{}
	)

	shortener.On("Get").Return("", errors.New("")).Once()
	handler := ShortenURL{
		shortener: shortener,
	}

	result := sendTestRequest(http.MethodGet, "/"+urlID, nil, handler.Get)
	assert.Equal(t, http.StatusNotFound, result.StatusCode)
	err := result.Body.Close()
	require.NoError(t, err)
	shortener.AssertExpectations(t)
}

func TestShortenURLHandler_isURL(t *testing.T) {
	shortener := ShortenURL{}

	tests := []struct {
		name string
		url  string
		want bool
	}{
		{
			name: "корректный URL",
			url:  "https://ya.ru/",
			want: true,
		},
		{
			name: "проверка присутствия схемы в URL",
			url:  "ya.ru",
			want: false,
		},
		{
			name: "проверка присутствия hostname в URL",
			url:  "/path/test",
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, shortener.isURL(tt.url))
		})
	}
}

func TestShortenURLHandler_prepareShortenURL(t *testing.T) {
	shortener := ShortenURL{}
	assert.Equal(t, "http://localhost/1", shortener.prepareShortenURL("localhost", "1"))
}

func sendTestRequest(method string, target string, body io.Reader, handler http.HandlerFunc) *http.Response {
	request := httptest.NewRequest(method, target, body)
	w := httptest.NewRecorder()
	handler(w, request)

	return w.Result()
}
