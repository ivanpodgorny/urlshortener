package storage

import (
	"context"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	inerr "github.com/ivanpodgorny/urlshortener/internal/app/errors"
)

func TestMemoryWithFile(t *testing.T) {
	var (
		filename          = "test"
		id                = "id1"
		wrongID           = "id2"
		idToDelete        = "id3"
		url               = "https://ya.ru/"
		userID            = "userID1"
		userWithoutURLsID = "userID2"
		ctx               = context.Background()
	)

	s, file := createFileStorage(t, filename)

	insertedID, err := s.Add(ctx, id, url, userID)
	assert.NoError(t, err, "добавление новой записи")
	assert.Equal(t, id, insertedID, "добавление новой записи")
	_, err = s.Add(ctx, id, url, userID)
	assert.Error(t, err, "добавление записи c существующим id")
	stored, err := s.Get(ctx, id)
	assert.NoError(t, err, "получение записи")
	assert.Equal(t, url, stored, "получение записи")
	_, err = s.Get(ctx, wrongID)
	assert.Error(t, err, "получение несуществующей записи")
	urls := s.GetAllUser(ctx, userID)
	assert.Equal(t, map[string]string{id: url}, urls, "получение URL пользователя")
	urls = s.GetAllUser(ctx, userWithoutURLsID)
	assert.Equal(t, map[string]string{}, urls, "получение URL пользователя, не добавлявшего URL")
	_, _ = s.Add(ctx, idToDelete, url, userID)
	err = s.DeleteBatch(ctx, []string{idToDelete}, userWithoutURLsID)
	assert.NoError(t, err, "попытка удаления чужой записи")
	notDeletedURL, err := s.Get(ctx, idToDelete)
	assert.NoError(t, err, "попытка удаления чужой записи")
	assert.Equal(t, url, notDeletedURL, "попытка удаления чужой записи")
	err = s.DeleteBatch(ctx, []string{idToDelete}, userID)
	assert.NoError(t, err, "удаление записи")
	_, err = s.Get(ctx, idToDelete)
	assert.ErrorIs(t, err, inerr.ErrURLIsDeleted, "получение удаленной записи")

	require.NoError(t, file.Close(), "не удалось закрыть файл")
	s, file = createFileStorage(t, filename)

	stored, err = s.Get(context.Background(), id)
	assert.NoError(t, err, "получение записи, сохраненной в файл")
	assert.Equal(t, url, stored, "получение записи, сохраненной в файл")
	urls = s.GetAllUser(context.Background(), userID)
	assert.Equal(t, map[string]string{id: url}, urls, "получение URL пользователя")
	_, err = s.Get(ctx, idToDelete)
	assert.ErrorIs(t, err, inerr.ErrURLIsDeleted, "получение удаленной записи")

	require.NoError(t, file.Close(), "не удалось закрыть файл")
	require.NoError(t, os.Remove(filename))
}

func TestMemoryOnly(t *testing.T) {
	var (
		id                = "id1"
		wrongID           = "id2"
		url               = "https://ya.ru/"
		userID            = "userID1"
		userWithoutURLsID = "userID2"
		ctx               = context.Background()
		s                 = NewMemory(nil)
	)

	insertedID, err := s.Add(ctx, id, url, userID)
	assert.NoError(t, err, "добавление новой записи")
	assert.Equal(t, id, insertedID, "добавление новой записи")
	_, err = s.Add(ctx, id, url, userID)
	assert.Error(t, err, "добавление записи c существующим id")
	stored, err := s.Get(ctx, id)
	assert.NoError(t, err, "получение записи")
	assert.Equal(t, url, stored, "получение записи")
	_, err = s.Get(ctx, wrongID)
	assert.Error(t, err, "получение несуществующей записи")
	urls := s.GetAllUser(ctx, userID)
	assert.Equal(t, map[string]string{id: url}, urls, "получение URL пользователя")
	urls = s.GetAllUser(ctx, userWithoutURLsID)
	assert.Equal(t, map[string]string{}, urls, "получение URL пользователя, не добавлявшего URL")
	err = s.DeleteBatch(ctx, []string{id}, userWithoutURLsID)
	assert.NoError(t, err, "попытка удаления чужой записи")
	notDeletedURL, err := s.Get(ctx, id)
	assert.NoError(t, err, "попытка удаления чужой записи")
	assert.Equal(t, url, notDeletedURL, "попытка удаления чужой записи")
	err = s.DeleteBatch(ctx, []string{id}, userID)
	assert.NoError(t, err, "удаление записи")
	_, err = s.Get(ctx, id)
	assert.ErrorIs(t, err, inerr.ErrURLIsDeleted, "получение удаленной записи")
}

func createFileStorage(t *testing.T, filename string) (*Memory, *os.File) {
	file, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0600)
	require.NoError(t, err, "не удалось создать файл")

	return NewMemory(file), file
}
