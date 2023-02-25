package service

import (
	"context"
	"errors"
	"github.com/ivanpodgorny/urlshortener/internal/app/storage"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

type StorageMock struct {
	mock.Mock
}

func (m *StorageMock) Add(_ context.Context, _ string, url string) error {
	args := m.Called(url)

	return args.Error(0)
}

func (m *StorageMock) Get(_ context.Context, _ string) (string, error) {
	return "", nil
}

func TestShortener_Shorten(t *testing.T) {
	shortener := Shortener{
		storage: storage.NewMemory(),
	}

	u := "https://ya.ru/"
	id, err := shortener.Shorten(context.Background(), u)
	assert.NoError(t, err)
	_, err = shortener.Get(context.Background(), "1")
	assert.Error(t, err)
	savedURL, err := shortener.Get(context.Background(), id)
	assert.NoError(t, err)
	assert.Equal(t, u, savedURL)
}

func TestShortener_ShortenReturnsError(t *testing.T) {
	var (
		noErrURL = "url1"
		errURL   = "url2"
	)
	storageMock := &StorageMock{}
	storageMock.
		On("Add", noErrURL).Return(nil).Once().
		On("Add", errURL).Return(errors.New("")).Once()
	shortener := Shortener{
		storage: storageMock,
	}

	_, err := shortener.Shorten(context.Background(), noErrURL)
	assert.NoError(t, err)
	_, err = shortener.Shorten(context.Background(), errURL)
	assert.Error(t, err)
}
