package main

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

func (s Shortener) Shorten(ctx context.Context, url string) (string, error) {
	id, err := generateId(12)
	if err != nil {
		return "", err
	}

	if err := s.storage.Add(ctx, id, url); err != nil {
		return "", err
	}

	return id, nil
}

func (s Shortener) Get(ctx context.Context, id string) (string, error) {
	return s.storage.Get(ctx, id)
}

func generateId(length int) (string, error) {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}

	return base64.URLEncoding.EncodeToString(b)[:length], nil
}
