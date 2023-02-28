package storage

import (
	"context"
)

// Memory реализует интерфейс service.Storage
// для хранения url в памяти в виде хэш-таблицы.
type Memory struct {
	urls     map[string]string
	userData map[string][]string
}

func NewMemory() *Memory {
	return &Memory{
		urls:     map[string]string{},
		userData: map[string][]string{},
	}
}

func (m Memory) Add(_ context.Context, id string, url string, userID string) error {
	if _, exist := m.urls[id]; exist {
		return ErrKeyExists
	}

	m.urls[id] = url
	m.userData[userID] = append(m.userData[userID], id)

	return nil
}

func (m Memory) Get(_ context.Context, id string) (string, error) {
	url, ok := m.urls[id]
	if !ok {
		return "", ErrKeyNotFound
	}

	return url, nil
}

func (m Memory) GetAllUser(ctx context.Context, userID string) map[string]string {
	data := map[string]string{}
	ids, ok := m.userData[userID]
	if !ok {
		return data
	}

	for _, id := range ids {
		url, err := m.Get(ctx, id)
		if err == nil {
			data[id] = url
		}
	}

	return data
}
