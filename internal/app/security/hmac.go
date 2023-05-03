package security

import (
	"crypto/hmac"
	"crypto/sha256"
)

// SignHMAC подписывает строку HMAC.
func SignHMAC(data []byte, key string) []byte {
	h := hmac.New(sha256.New, []byte(key))
	h.Write(data)

	return h.Sum(nil)
}

// ValidateHMAC проверяет подлинность подписи HMAC.
func ValidateHMAC(data, sign []byte, key string) bool {
	return hmac.Equal(SignHMAC(data, key), sign)
}
