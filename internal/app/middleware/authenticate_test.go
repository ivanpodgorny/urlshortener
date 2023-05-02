package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type AuthenticatorMock struct {
	mock.Mock
}

func (m *AuthenticatorMock) Authenticate(w http.ResponseWriter, r *http.Request) *http.Request {
	m.Called(w, r)

	return r
}

func TestAuthenticate(t *testing.T) {
	var (
		r             = chi.NewRouter()
		path          = "/"
		authenticator = &AuthenticatorMock{}
		successStatus = http.StatusOK
	)

	authenticator.
		On("Authenticate", mock.AnythingOfType("*http.response"), mock.AnythingOfType("*http.Request")).
		Return(mock.AnythingOfType("*http.Request")).
		Once()
	r.Use(Authenticate(authenticator))
	r.Post(path, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(successStatus)
	})

	ts := httptest.NewServer(r)
	defer ts.Close()

	req, err := http.NewRequest(http.MethodPost, ts.URL+path, nil)
	require.NoError(t, err)
	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	require.NoError(t, resp.Body.Close())
	assert.Equal(t, resp.StatusCode, successStatus)
	authenticator.AssertExpectations(t)
}
