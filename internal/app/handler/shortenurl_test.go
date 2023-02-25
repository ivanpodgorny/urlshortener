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

func (m *ShortenerMock) Shorten(_ context.Context, _, userID string) (string, error) {
	args := m.Called(userID)

	return args.String(0), args.Error(1)
}

func (m *ShortenerMock) Get(_ context.Context, _ string) (string, error) {
	args := m.Called()

	return args.String(0), args.Error(1)
}

func (m *ShortenerMock) GetAllUser(_ context.Context, userID string) map[string]string {
	args := m.Called(userID)

	return args.Get(0).(map[string]string)
}

type AuthenticatorMock struct {
	mock.Mock
}

func (m *AuthenticatorMock) UserIdentifier(_ *http.Request) (string, error) {
	args := m.Called()

	return args.String(0), args.Error(1)
}

func TestShortenURLHandler_CreateSuccess(t *testing.T) {
	var (
		urlID         = "1i-CBrzwyMkL"
		userID        = "438c4b98-fc98-45cf-ac63-c4a86fbd4ff4"
		baseURL       = "http://localhost"
		shortener     = &ShortenerMock{}
		authenticator = &AuthenticatorMock{}
	)

	authenticator.On("UserIdentifier").Return(userID, nil).Once()
	shortener.On("Shorten", userID).Return(urlID, nil).Once()
	handler := ShortenURL{
		shortener:     shortener,
		baseURL:       baseURL,
		authenticator: authenticator,
	}

	result := sendTestRequest(http.MethodPost, "/", bytes.NewBuffer([]byte("https://ya.ru/")), handler.Create)
	assert.Equal(t, http.StatusCreated, result.StatusCode)
	b, err := io.ReadAll(result.Body)
	require.NoError(t, err)
	assert.Equal(t, baseURL+"/"+urlID, string(b[:]))
	err = result.Body.Close()
	require.NoError(t, err)
	authenticator.AssertExpectations(t)
	shortener.AssertExpectations(t)
}

func TestShortenURLHandler_CreateWithErrors(t *testing.T) {
	var (
		userID        = "438c4b98-fc98-45cf-ac63-c4a86fbd4ff4"
		shortener     = &ShortenerMock{}
		authenticator = &AuthenticatorMock{}
	)

	authenticator.On("UserIdentifier").Return(userID, nil).Times(3)
	shortener.On("Shorten", userID).Return("", errors.New("")).Once()
	handler := ShortenURL{
		shortener:     shortener,
		authenticator: authenticator,
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
	authenticator.AssertExpectations(t)
	shortener.AssertExpectations(t)
}

func TestShortenURLHandler_CreateJSONSuccess(t *testing.T) {
	var (
		urlID         = "1i-CBrzwyMkL"
		userID        = "438c4b98-fc98-45cf-ac63-c4a86fbd4ff4"
		baseURL       = "http://localhost"
		shortener     = &ShortenerMock{}
		authenticator = &AuthenticatorMock{}
	)

	authenticator.On("UserIdentifier").Return(userID, nil).Once()
	shortener.On("Shorten", userID).Return(urlID, nil).Once()
	handler := ShortenURL{
		shortener:     shortener,
		baseURL:       baseURL,
		authenticator: authenticator,
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
	assert.Equal(t, baseURL+"/"+urlID, resp.URL)
	err = result.Body.Close()
	require.NoError(t, err)
	authenticator.AssertExpectations(t)
	shortener.AssertExpectations(t)
}

func TestShortenURLHandler_CreateJSONWithErrors(t *testing.T) {
	var (
		userID        = "438c4b98-fc98-45cf-ac63-c4a86fbd4ff4"
		shortener     = &ShortenerMock{}
		authenticator = &AuthenticatorMock{}
	)

	authenticator.On("UserIdentifier").Return(userID, nil).Times(3)
	shortener.On("Shorten", userID).Return("", errors.New("")).Once()
	handler := ShortenURL{
		shortener:     shortener,
		authenticator: authenticator,
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
	authenticator.AssertExpectations(t)
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

func TestShortenURL_GetAllByCurrentUserSuccess(t *testing.T) {
	var (
		urlID         = "1i-CBrzwyMkL"
		url           = "https://ya.ru/"
		urls          = map[string]string{urlID: url}
		userID        = "438c4b98-fc98-45cf-ac63-c4a86fbd4ff4"
		baseURL       = "http://localhost"
		shortener     = &ShortenerMock{}
		authenticator = &AuthenticatorMock{}
	)

	authenticator.On("UserIdentifier").Return(userID, nil).Once()
	shortener.On("GetAllUser", userID).Return(urls, nil).Once()
	handler := ShortenURL{
		shortener:     shortener,
		baseURL:       baseURL,
		authenticator: authenticator,
	}

	result := sendTestRequest(http.MethodPost, "/", nil, handler.GetAllByCurrentUser)
	assert.Equal(t, http.StatusOK, result.StatusCode)
	assert.Equal(t, "application/json", result.Header.Get("Content-Type"))
	b, err := io.ReadAll(result.Body)
	require.NoError(t, err)
	resp := make([]struct {
		ShortURL    string `json:"short_url"`
		OriginalURL string `json:"original_url"`
	}, 0)
	err = json.Unmarshal(b, &resp)
	require.NoError(t, err)
	urlData := resp[0]
	assert.Equal(t, baseURL+"/"+urlID, urlData.ShortURL)
	assert.Equal(t, url, urlData.OriginalURL)
	err = result.Body.Close()
	require.NoError(t, err)
	authenticator.AssertExpectations(t)
	shortener.AssertExpectations(t)
}

func TestShortenURL_GetAllByCurrentUserNoContent(t *testing.T) {
	var (
		urls          = map[string]string{}
		userID        = "438c4b98-fc98-45cf-ac63-c4a86fbd4ff4"
		shortener     = &ShortenerMock{}
		authenticator = &AuthenticatorMock{}
	)

	authenticator.On("UserIdentifier").Return(userID, nil).Once()
	shortener.On("GetAllUser", userID).Return(urls, nil).Once()
	handler := ShortenURL{
		shortener:     shortener,
		authenticator: authenticator,
	}

	result := sendTestRequest(http.MethodPost, "/", nil, handler.GetAllByCurrentUser)
	assert.Equal(t, http.StatusNoContent, result.StatusCode)
	err := result.Body.Close()
	require.NoError(t, err)
	authenticator.AssertExpectations(t)
	shortener.AssertExpectations(t)
}

func TestUserAuthenticationErrors(t *testing.T) {
	authenticator := &AuthenticatorMock{}
	authenticator.On("UserIdentifier").Return("userID", errors.New("")).Times(3)
	handler := ShortenURL{
		authenticator: authenticator,
	}

	tests := []struct {
		name    string
		handler http.HandlerFunc
	}{
		{
			name:    "Create",
			handler: handler.Create,
		},
		{
			name:    "CreateJSON",
			handler: handler.CreateJSON,
		},
		{
			name:    "GetAllByCurrentUser",
			handler: handler.GetAllByCurrentUser,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := sendTestRequest(http.MethodPost, "/", nil, tt.handler)
			assert.Equal(t, http.StatusUnauthorized, result.StatusCode)
			err := result.Body.Close()
			require.NoError(t, err)
		})
	}
	authenticator.AssertExpectations(t)
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
	shortener := ShortenURL{
		baseURL: "http://localhost",
	}
	assert.Equal(t, "http://localhost/1", shortener.prepareShortenURL("1"))
}

func sendTestRequest(method string, target string, body io.Reader, handler http.HandlerFunc) *http.Response {
	request := httptest.NewRequest(method, target, body)
	w := httptest.NewRecorder()
	handler(w, request)

	return w.Result()
}
