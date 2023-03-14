package service

import (
	"context"
	"github.com/ivanpodgorny/urlshortener/internal/app/security"
)

type Shortener struct {
	storage Storage
}

type Storage interface {
	Add(ctx context.Context, id string, url string, userID string) (string, error)
	Get(ctx context.Context, id string) (string, error)
	GetAllUser(ctx context.Context, userID string) map[string]string
	DeleteBatch(ctx context.Context, urlIDs []string, userID string) error
}

func NewShortener(s Storage) *Shortener {
	return &Shortener{storage: s}
}

// Shorten принимает строку URL, генерирует для нее случайный текстовый ID,
// сохраняет ID и URL в Storage и возвращает сгенерированный ID.
// Если URL уже сохранен в Storage, новая запись не добавляется и во втором параметре вернется false.
// Если сгенерированный ID уже существует в Storage, возвращает ошибку.
func (s Shortener) Shorten(ctx context.Context, url string, userID string) (string, bool, error) {
	id, err := security.GenerateRandomString(16)
	if err != nil {
		return "", false, err
	}

	storedID, err := s.storage.Add(ctx, id, url, userID)
	if err != nil {
		return "", false, err
	}

	return storedID, storedID == id, nil
}

// Get принимает текстовый ID и возвращает URL, сохраненный в Storage с этим ID.
func (s Shortener) Get(ctx context.Context, id string) (string, error) {
	return s.storage.Get(ctx, id)
}

// GetAllUser принимает идентификатор пользователя
// и возвращает сокращенные им URL и их ID в формате {ID: URL, ...}.
func (s Shortener) GetAllUser(ctx context.Context, userID string) map[string]string {
	return s.storage.GetAllUser(ctx, userID)
}

// DeleteBatch принимает массив идентификаторов URL и выполняет их удаление из Storage.
func (s Shortener) DeleteBatch(ctx context.Context, urlIDs []string, userID string) error {
	return s.storage.DeleteBatch(ctx, urlIDs, userID)
}
