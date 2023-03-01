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
	"testing"
)

type ShortenerMock struct {
	mock.Mock
}

func (m *ShortenerMock) Shorten(_ context.Context, url, userID string) (string, bool, error) {
	args := m.Called(url, userID)

	return args.String(0), args.Bool(1), args.Error(2)
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
		url           = "https://ya.ru/"
		userID        = "438c4b98-fc98-45cf-ac63-c4a86fbd4ff4"
		baseURL       = "http://localhost"
		shortener     = &ShortenerMock{}
		authenticator = &AuthenticatorMock{}
	)

	authenticator.On("UserIdentifier").Return(userID, nil).Once()
	shortener.On("Shorten", url, userID).Return(urlID, true, nil).Once()
	handler := ShortenURL{
		shortener:     shortener,
		baseURL:       baseURL,
		authenticator: authenticator,
	}

	result := sendTestRequest(http.MethodPost, "/", bytes.NewBuffer([]byte(url)), handler.Create)
	assert.Equal(t, http.StatusCreated, result.StatusCode)
	b, err := io.ReadAll(result.Body)
	require.NoError(t, err)
	assert.Equal(t, baseURL+"/"+urlID, string(b[:]))
	err = result.Body.Close()
	require.NoError(t, err)
	authenticator.AssertExpectations(t)
	shortener.AssertExpectations(t)
}

func TestShortenURLHandler_CreateConflict(t *testing.T) {
	var (
		url           = "https://ya.ru/"
		userID        = "438c4b98-fc98-45cf-ac63-c4a86fbd4ff4"
		shortener     = &ShortenerMock{}
		authenticator = &AuthenticatorMock{}
	)

	authenticator.On("UserIdentifier").Return(userID, nil).Once()
	shortener.On("Shorten", url, userID).Return("", false, nil).Once()
	handler := ShortenURL{
		shortener:     shortener,
		authenticator: authenticator,
	}

	result := sendTestRequest(http.MethodPost, "/", bytes.NewBuffer([]byte(url)), handler.Create)
	assert.Equal(t, http.StatusConflict, result.StatusCode)
	require.NoError(t, result.Body.Close())
	authenticator.AssertExpectations(t)
	shortener.AssertExpectations(t)
}

func TestShortenURLHandler_CreateWithErrors(t *testing.T) {
	var (
		userID        = "438c4b98-fc98-45cf-ac63-c4a86fbd4ff4"
		url           = "https://ya.ru/"
		shortener     = &ShortenerMock{}
		authenticator = &AuthenticatorMock{}
	)

	authenticator.On("UserIdentifier").Return(userID, nil).Times(3)
	shortener.On("Shorten", url, userID).Return("", false, errors.New("")).Once()
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
			body:           bytes.NewBuffer([]byte(url)),
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
		url           = "https://ya.ru/"
		userID        = "438c4b98-fc98-45cf-ac63-c4a86fbd4ff4"
		baseURL       = "http://localhost"
		shortener     = &ShortenerMock{}
		authenticator = &AuthenticatorMock{}
	)

	authenticator.On("UserIdentifier").Return(userID, nil).Once()
	shortener.On("Shorten", url, userID).Return(urlID, true, nil).Once()
	handler := ShortenURL{
		shortener:     shortener,
		baseURL:       baseURL,
		authenticator: authenticator,
	}

	result := sendTestRequest(http.MethodPost, "/", bytes.NewBuffer([]byte(`{"url":"`+url+`"}`)), handler.CreateJSON)
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

func TestShortenURLHandler_CreateJSONConflict(t *testing.T) {
	var (
		url           = "https://ya.ru/"
		userID        = "438c4b98-fc98-45cf-ac63-c4a86fbd4ff4"
		shortener     = &ShortenerMock{}
		authenticator = &AuthenticatorMock{}
	)

	authenticator.On("UserIdentifier").Return(userID, nil).Once()
	shortener.On("Shorten", url, userID).Return("", false, nil).Once()
	handler := ShortenURL{
		shortener:     shortener,
		authenticator: authenticator,
	}

	result := sendTestRequest(http.MethodPost, "/", bytes.NewBuffer([]byte(`{"url":"`+url+`"}`)), handler.CreateJSON)
	assert.Equal(t, http.StatusConflict, result.StatusCode)
	assert.Equal(t, "application/json", result.Header.Get("Content-Type"))
	require.NoError(t, result.Body.Close())
	authenticator.AssertExpectations(t)
	shortener.AssertExpectations(t)
}

func TestShortenURLHandler_CreateJSONWithErrors(t *testing.T) {
	var (
		userID        = "438c4b98-fc98-45cf-ac63-c4a86fbd4ff4"
		url           = "https://ya.ru/"
		shortener     = &ShortenerMock{}
		authenticator = &AuthenticatorMock{}
	)

	authenticator.On("UserIdentifier").Return(userID, nil).Times(3)
	shortener.On("Shorten", url, userID).Return("", false, errors.New("")).Once()
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
			body:           bytes.NewBuffer([]byte(`{"url":"` + url + `"}`)),
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

func TestShortenURLHandler_CreateBatchSuccess(t *testing.T) {
	var (
		id1           = "1"
		id2           = "2"
		urlID1        = "InFsfCVTdXY7cVly"
		urlID2        = "SVp4LEsxaE1EWVMq"
		url1          = "https://ya.ru/"
		url2          = "https://www.google.ru/"
		userID        = "438c4b98-fc98-45cf-ac63-c4a86fbd4ff4"
		baseURL       = "http://localhost"
		bodyJSON      = `[{"correlation_id":"` + id1 + `","original_url":"` + url1 + `"},{"correlation_id": "` + id2 + `","original_url":"` + url2 + `"}]`
		shortener     = &ShortenerMock{}
		authenticator = &AuthenticatorMock{}
	)

	authenticator.On("UserIdentifier").Return(userID, nil).Once()
	shortener.On("Shorten", url1, userID).Return(urlID1, true, nil).Once()
	shortener.On("Shorten", url2, userID).Return(urlID2, true, nil).Once()
	handler := ShortenURL{
		shortener:     shortener,
		baseURL:       baseURL,
		authenticator: authenticator,
	}

	result := sendTestRequest(http.MethodPost, "/", bytes.NewBuffer([]byte(bodyJSON)), handler.CreateBatch)
	assert.Equal(t, http.StatusCreated, result.StatusCode)
	assert.Equal(t, "application/json", result.Header.Get("Content-Type"))
	b, err := io.ReadAll(result.Body)
	require.NoError(t, err)
	resp := make([]struct {
		ID  string `json:"correlation_id"`
		URL string `json:"short_url"`
	}, 0)
	err = json.Unmarshal(b, &resp)
	require.NoError(t, err)

	foundID1, foundID2 := false, false
	for _, u := range resp {
		if u.ID == id1 {
			assert.Equal(t, baseURL+"/"+urlID1, u.URL, "URL не соотвествует идентификатору")
			foundID1 = true

			continue
		}
		if u.ID == id2 {
			assert.Equal(t, baseURL+"/"+urlID2, u.URL, "URL не соотвествует идентификатору")
			foundID2 = true
		}
	}
	assert.True(t, foundID1, "не найден идентификатор оригинального URL")
	assert.True(t, foundID2, "не найден идентификатор оригинального URL")

	err = result.Body.Close()
	require.NoError(t, err)
	authenticator.AssertExpectations(t)
	shortener.AssertExpectations(t)
}

func TestShortenURLHandler_CreateBatchError(t *testing.T) {
	var (
		userID        = "438c4b98-fc98-45cf-ac63-c4a86fbd4ff4"
		baseURL       = "http://localhost"
		authenticator = &AuthenticatorMock{}
	)

	authenticator.On("UserIdentifier").Return(userID, nil).Once()
	handler := ShortenURL{
		baseURL:       baseURL,
		authenticator: authenticator,
	}

	result := sendTestRequest(http.MethodPost, "/", bytes.NewBuffer([]byte("{}")), handler.CreateBatch)
	assert.Equal(t, http.StatusBadRequest, result.StatusCode)
	require.NoError(t, result.Body.Close())
	authenticator.AssertExpectations(t)
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
			assert.Equal(t, tt.want, shortener.validateURL(tt.url))
		})
	}
}

func TestShortenURLHandler_prepareShortenURL(t *testing.T) {
	shortener := ShortenURL{
		baseURL: "http://localhost",
	}
	assert.Equal(t, "http://localhost/1", shortener.prepareShortenURL("1"))
}
