package service

import (
	"context"
	"github.com/ivanpodgorny/urlshortener/internal/app/security"
)

type Shortener struct {
	storage Storage
}

type Storage interface {
	Add(ctx context.Context, id string, url string, userID string) error
	Get(ctx context.Context, id string) (string, error)
	GetAllUser(ctx context.Context, userID string) map[string]string
}

func NewShortener(s Storage) *Shortener {
	return &Shortener{storage: s}
}

// Shorten принимает строку URL, генерирует для нее случайный текстовый ID,
// сохраняет ID и URL в Storage и возвращает сгенерированный ID.
// Если сгенерированный ID уже существует в Storage, возвращает ошибку.
func (s Shortener) Shorten(ctx context.Context, url string, userID string) (string, error) {
	id, err := security.GenerateRandomString(16)
	if err != nil {
		return "", err
	}

	if err := s.storage.Add(ctx, id, url, userID); err != nil {
		return "", err
	}

	return id, nil
}

// Get принимает текстовый ID и возвращает URL, сохраненный в Storage с этим ID.
func (s Shortener) Get(ctx context.Context, id string) (string, error) {
	return s.storage.Get(ctx, id)
}

func (s Shortener) GetAllUser(ctx context.Context, userID string) map[string]string {
	return s.storage.GetAllUser(ctx, userID)
}
