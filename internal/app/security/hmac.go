package security

import (
	"crypto/hmac"
	"crypto/sha256"
)

func SignHMAC(data []byte, key string) []byte {
	h := hmac.New(sha256.New, []byte(key))
	h.Write(data)

	return h.Sum(nil)
}

func ValidateHMAC(data, sign []byte, key string) bool {
	return hmac.Equal(SignHMAC(data, key), sign)
}
