package middleware

import (
	"net"
	"net/http"
)

// Internal возвращает middleware, проверяющий, что переданный в заголовке запроса
// X-Real-IP IP-адрес клиента входит в доверенную подсеть trustedSubnet.
func Internal(trustedSubnet string) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ip := net.ParseIP(r.Header.Get("X-Real-IP"))
			_, ipNet, err := net.ParseCIDR(trustedSubnet)
			if ip == nil || err != nil || !ipNet.Contains(ip) {
				w.WriteHeader(http.StatusForbidden)

				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
