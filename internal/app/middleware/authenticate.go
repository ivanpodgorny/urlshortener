package middleware

import (
	"net/http"
)

// Authenticator интерфейс сервиса аутентификации пользователя.
type Authenticator interface {
	Authenticate(http.ResponseWriter, *http.Request) *http.Request
}

// Authenticate возвращает middleware для поверки токена пользователя.
func Authenticate(a Authenticator) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			r = a.Authenticate(w, r)
			next.ServeHTTP(w, r)
		})
	}
}
