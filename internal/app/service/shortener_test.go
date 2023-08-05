package service

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	inerr "github.com/ivanpodgorny/urlshortener/internal/app/errors"
)

type StorageMock struct {
	mock.Mock
}

func (m *StorageMock) Add(_ context.Context, id string, url string, userID string) (string, error) {
	args := m.Called(url, userID)

	return id, args.Error(0)
}

func (m *StorageMock) Get(_ context.Context, id string) (string, error) {
	args := m.Called(id)

	return args.String(0), args.Error(1)
}

func (m *StorageMock) GetAllUser(_ context.Context, userID string) map[string]string {
	args := m.Called(userID)

	return args.Get(0).(map[string]string)
}

func (m *StorageMock) DeleteBatch(_ context.Context, urlIDs []string, userID string) error {
	args := m.Called(urlIDs, userID)

	return args.Error(0)
}

func (m *StorageMock) GetStat(_ context.Context) (int, int, error) {
	args := m.Called()

	return args.Int(0), args.Int(1), args.Error(2)
}

func TestShortener(t *testing.T) {
	var (
		userID     = "438c4b98-fc98-45cf-ac63-c4a86fbd4ff4"
		url        = "https://ya.ru/"
		urlID      = "1i-CBrzwyMkL"
		urlIDs     = []string{urlID}
		urls       = map[string]string{urlID: url}
		urlCount   = 2
		usersCount = 1
		ctx        = context.Background()
		storage    = &StorageMock{}
	)

	storage.
		On("Add", url, userID).Return(nil).Once().
		On("Get", urlID).Return(url, nil).Once().
		On("GetAllUser", userID).Return(urls).Once().
		On("DeleteBatch", urlIDs, userID).Return(nil).Once().
		On("GetStat").Return(urlCount, usersCount, nil).Once()
	shortener := Shortener{
		storage: storage,
	}

	_, inserted, err := shortener.Shorten(ctx, url, userID)
	assert.NoError(t, err)
	assert.True(t, inserted)
	savedURL, err := shortener.Get(ctx, urlID)
	assert.NoError(t, err)
	assert.Equal(t, url, savedURL)
	userURLs := shortener.GetAllUser(ctx, userID)
	assert.Equal(t, urls, userURLs)
	err = shortener.DeleteBatch(ctx, urlIDs, userID)
	assert.NoError(t, err)
	getURLCount, getUsersCount, err := shortener.GetStat(ctx)
	assert.NoError(t, err)
	assert.Equal(t, urlCount, getURLCount)
	assert.Equal(t, usersCount, getUsersCount)
	storage.AssertExpectations(t)
}

func TestShortenerReturnsError(t *testing.T) {
	var (
		url     = "https://ya.ru/"
		urlID   = "1i-CBrzwyMkL"
		urlIDs  = []string{urlID}
		userID  = "438c4b98-fc98-45cf-ac63-c4a86fbd4ff4"
		ctx     = context.Background()
		storage = &StorageMock{}
	)

	storage.
		On("Add", url, userID).Return(errors.New("")).Once().
		On("Get", urlID).Return("", errors.New("")).Once().
		On("DeleteBatch", urlIDs, userID).Return(inerr.ErrURLIsDeleted).Once().
		On("GetStat").Return(0, 0, errors.New("")).Once()
	shortener := Shortener{
		storage: storage,
	}

	_, _, err := shortener.Shorten(ctx, url, userID)
	assert.Error(t, err)
	_, err = shortener.Get(ctx, urlID)
	assert.Error(t, err)
	err = shortener.DeleteBatch(ctx, urlIDs, userID)
	assert.ErrorIs(t, err, inerr.ErrURLIsDeleted)
	_, _, err = shortener.GetStat(ctx)
	assert.Error(t, err)
	storage.AssertExpectations(t)
}
