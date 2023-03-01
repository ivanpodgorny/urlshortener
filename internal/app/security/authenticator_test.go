package security

import (
	"github.com/stretchr/testify/assert"
	"net/http/httptest"
	"testing"
)

func TestAuthenticator(t *testing.T) {
	var (
		authenticator = NewAuthenticator(
			NewCookieTokenStorage(NewHMACTokenCreatorParser("")),
			RequestContextUserProvider{},
		)
		request  = httptest.NewRequest("", "/", nil)
		recorder = httptest.NewRecorder()
	)

	_, err := authenticator.UserIdentifier(request)
	assert.Error(t, err, "неаутентифицированный пользователь")

	request = authenticator.Authenticate(recorder, request)
	id, err := authenticator.UserIdentifier(request)
	assert.NoError(t, err, "создание нового токена для пользователя")

	resp := recorder.Result()
	defer resp.Body.Close()

	request.AddCookie(resp.Cookies()[0])
	request = authenticator.Authenticate(recorder, request)
	idFromCookies, err := authenticator.UserIdentifier(request)
	assert.NoError(t, err, "получение существующего токена пользователя")
	assert.Equal(t, id, idFromCookies, "получение существующего токена пользователя")
}
