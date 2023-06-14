package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInternal(t *testing.T) {
	tests := []struct {
		name           string
		trustedSubnet  string
		wantStatusCode int
	}{
		{
			name:           "пустое значение trustedSubnet",
			trustedSubnet:  "",
			wantStatusCode: http.StatusForbidden,
		},
		{
			name:           "адрес входит в подсеть",
			trustedSubnet:  "127.0.0.0/24",
			wantStatusCode: http.StatusOK,
		},
		{
			name:           "адрес не входит в подсеть",
			trustedSubnet:  "192.168.0.0/24",
			wantStatusCode: http.StatusForbidden,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sendInternalTestRequest(t, tt.trustedSubnet, tt.wantStatusCode)
		})
	}
}

func sendInternalTestRequest(t *testing.T, trustedSubnet string, wantStatusCode int) {
	var (
		r    = chi.NewRouter()
		path = "/"
	)

	r.Use(Internal(trustedSubnet))
	r.Get(path, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	ts := httptest.NewServer(r)
	defer ts.Close()

	req, err := http.NewRequest(http.MethodGet, ts.URL+path, nil)
	require.NoError(t, err)
	req.Header.Set("X-Real-IP", "127.0.0.1")
	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	require.NoError(t, resp.Body.Close())
	assert.Equal(t, wantStatusCode, resp.StatusCode)
}
