package security

import (
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/stretchr/testify/assert"
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

	_, err := authenticator.UserIdentifier(request.Context())
	assert.Error(t, err, "неаутентифицированный пользователь")

	request = authenticator.Authenticate(recorder, request)
	id, err := authenticator.UserIdentifier(request.Context())
	assert.NoError(t, err, "создание нового токена для пользователя")

	resp := recorder.Result()
	defer require.NoError(t, resp.Body.Close())

	request.AddCookie(resp.Cookies()[0])
	request = authenticator.Authenticate(recorder, request)
	idFromCookies, err := authenticator.UserIdentifier(request.Context())
	assert.NoError(t, err, "получение существующего токена пользователя")
	assert.Equal(t, id, idFromCookies, "получение существующего токена пользователя")
}
