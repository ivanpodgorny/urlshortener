package storage

import (
	"context"
)

type Memory struct {
	data map[string]string
}

func NewMemory() *Memory {
	return &Memory{
		data: map[string]string{},
	}
}

func (m Memory) Add(_ context.Context, id string, url string) error {
	if _, exist := m.data[id]; exist {
		return ErrKeyExists
	}

	m.data[id] = url

	return nil
}

func (m Memory) Get(_ context.Context, id string) (string, error) {
	url, ok := m.data[id]
	if !ok {
		return "", ErrKeyNotFound
	}

	return url, nil
}
