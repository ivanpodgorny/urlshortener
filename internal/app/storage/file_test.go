package storage

import (
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"os"
	"testing"
)

func TestFile(t *testing.T) {
	var (
		filename          = "test"
		id                = "id1"
		wrongID           = "id2"
		url               = "https://ya.ru/"
		userID            = "userID1"
		userWithoutURLsID = "userID2"
	)

	s, err := NewFile(filename)
	require.NoError(t, err, "не удалось создать файл")
	assert.NoError(t, s.Add(context.Background(), id, url, userID), "добавление новой записи")
	assert.Error(t, s.Add(context.Background(), id, url, userID), "добавление записи c существующим id")
	stored, err := s.Get(context.Background(), id)
	assert.NoError(t, err, "получение записи")
	assert.Equal(t, url, stored, "получение записи")
	_, err = s.Get(context.Background(), wrongID)
	assert.Error(t, err, "получение несуществующей записи")
	urls := s.GetAllUser(context.Background(), userID)
	assert.Equal(t, map[string]string{id: url}, urls, "получение URL пользователя")
	urls = s.GetAllUser(context.Background(), userWithoutURLsID)
	assert.Equal(t, map[string]string{}, urls, "получение URL пользователя, не добавлявшего URL")

	require.NoError(t, s.Close(), "не удалось закрыть файл")
	s, err = NewFile(filename)
	require.NoError(t, err, "не удалось прочитать файл")

	stored, err = s.Get(context.Background(), id)
	assert.NoError(t, err, "получение записи, сохраненной в файл")
	assert.Equal(t, url, stored, "получение записи, сохраненной в файл")
	urls = s.GetAllUser(context.Background(), userID)
	assert.Equal(t, map[string]string{id: url}, urls, "получение URL пользователя")

	require.NoError(t, s.Close(), "не удалось закрыть файл")
	require.NoError(t, os.Remove(filename))
}
