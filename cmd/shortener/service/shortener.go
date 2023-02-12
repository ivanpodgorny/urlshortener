package service

import (
	"context"
	"crypto/rand"
	"encoding/base64"
)

type Shortener struct {
	storage Storage
}

type Storage interface {
	Add(ctx context.Context, id string, url string) error
	Get(ctx context.Context, id string) (string, error)
}

func NewShortener(s Storage) *Shortener {
	return &Shortener{storage: s}
}

// Shorten принимает строку URL, генерирует для нее случайный текстовый ID,
// сохраняет ID и URL в Storage и возвращает сгенерированный ID.
// Если сгенерированный ID уже существует в Storage, возвращает ошибку.
func (s Shortener) Shorten(ctx context.Context, url string) (string, error) {
	id, err := generateID(12)
	if err != nil {
		return "", err
	}

	if err := s.storage.Add(ctx, id, url); err != nil {
		return "", err
	}

	return id, nil
}

// Get принимает текстовый ID и возвращает URL, сохраненный в Storage с этим ID.
func (s Shortener) Get(ctx context.Context, id string) (string, error) {
	return s.storage.Get(ctx, id)
}

func generateID(length int) (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}

	return base64.URLEncoding.EncodeToString(b)[:length], nil
}
