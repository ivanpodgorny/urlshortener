package security

import (
	"crypto/rand"
	"encoding/base64"
	"github.com/google/uuid"
)

func GenerateUUID() string {
	return uuid.New().String()
}

func GenerateRandomBytes(size int) ([]byte, error) {
	b := make([]byte, size)
	if _, err := rand.Read(b); err != nil {
		return nil, err
	}

	return b, nil
}

func GenerateRandomString(size int) (string, error) {
	b, err := GenerateRandomBytes(size)
	if err != nil {
		return "", err
	}

	return base64.RawURLEncoding.EncodeToString(b)[:], nil
}
