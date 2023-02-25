package middleware

import (
	"bytes"
	"compress/gzip"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestDecompress(t *testing.T) {
	var (
		url           = "https://ya.ru/"
		r             = chi.NewRouter()
		path          = "/"
		successStatus = http.StatusCreated
	)

	r.Use(Decompress())

	r.Post(path, func(w http.ResponseWriter, r *http.Request) {
		b, err := io.ReadAll(r.Body)
		require.NoError(t, err)
		assert.Equal(t, url, string(b), "чтение тела запроса с Content-Encoding: "+r.Header.Get("Content-Encoding"))
		w.WriteHeader(successStatus)
	})

	ts := httptest.NewServer(r)
	defer ts.Close()

	assert.Equal(
		t,
		sendTestRequest(t, ts, path, url, "gzip"),
		successStatus,
		"обработка Content-Encoding: gzip",
	)
	assert.Equal(
		t,
		sendTestRequest(t, ts, path, url, "identity"),
		successStatus,
		"обработка Content-Encoding: identity",
	)
	assert.Equal(
		t,
		sendTestRequest(t, ts, path, url, ""),
		successStatus,
		"обработка запроса без Content-Encoding",
	)
	assert.Equal(
		t,
		sendTestRequest(t, ts, path, url, "br"),
		http.StatusUnsupportedMediaType,
		"обработка запроса с неподдерживаемым Content-Encoding",
	)
}

func sendTestRequest(t *testing.T, ts *httptest.Server, path string, body string, enc string) int {
	req, err := http.NewRequest(http.MethodPost, ts.URL+path, prepareBody(t, body, enc))
	require.NoError(t, err)
	if enc != "" {
		req.Header.Set("Content-Encoding", enc)
	}

	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	require.NoError(t, resp.Body.Close())

	return resp.StatusCode
}

func prepareBody(t *testing.T, body string, enc string) io.Reader {
	var buf bytes.Buffer
	if enc == "identity" || enc == "" {
		buf.Write([]byte(body))

		return &buf
	}

	gz := gzip.NewWriter(&buf)
	_, err := gz.Write([]byte(body))
	require.NoError(t, err)
	require.NoError(t, gz.Close())

	return &buf
}
