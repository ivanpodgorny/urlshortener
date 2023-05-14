package security

import (
	"crypto/rand"
	"encoding/base64"

	"github.com/google/uuid"
)

// GenerateUUID генерирует UUID v4.
func GenerateUUID() string {
	return uuid.New().String()
}

// GenerateRandomBytes генерирует случайную последовательность байт длины size.
func GenerateRandomBytes(size int) ([]byte, error) {
	b := make([]byte, size)
	if _, err := rand.Read(b); err != nil {
		return nil, err
	}

	return b, nil
}

// GenerateRandomString генерирует случайную строку из цифр, букв латинского алфавита
// и символов "-_" длины size.
func GenerateRandomString(size int) (string, error) {
	b, err := GenerateRandomBytes(size)
	if err != nil {
		return "", err
	}

	return base64.RawURLEncoding.EncodeToString(b)[:size], nil
}
