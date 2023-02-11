package storage

import (
	"context"
	"errors"
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
	_, exist := m.data[id]
	if exist {
		return errors.New("key exists")
	}

	m.data[id] = url

	return nil
}

func (m Memory) Get(_ context.Context, id string) (string, error) {
	url, ok := m.data[id]
	if !ok {
		return "", errors.New("key not found")
	}

	return url, nil
}
