package handler

import (
	"io"
	"net/http"
	"net/http/httptest"
)

func sendTestRequest(method string, target string, body io.Reader, handler http.HandlerFunc) *http.Response {
	request := httptest.NewRequest(method, target, body)
	w := httptest.NewRecorder()
	handler(w, request)

	return w.Result()
}
