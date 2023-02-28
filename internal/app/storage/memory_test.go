package storage

import (
	"context"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMemory(t *testing.T) {
	var (
		id                = "id1"
		wrongID           = "id2"
		url               = "https://ya.ru/"
		userID            = "userID1"
		userWithoutURLsID = "userID2"
		s                 = NewMemory()
	)

	assert.NoError(t, s.Add(context.Background(), id, url, userID), "добавление новой записи")
	assert.Error(t, s.Add(context.Background(), id, url, userID), "добавление записи c существующим id")
	stored, err := s.Get(context.Background(), id)
	assert.NoError(t, err, "получение записи")
	assert.Equal(t, url, stored, "получение записи")
	_, err = s.Get(context.Background(), wrongID)
	assert.Error(t, err, "получение несуществующей записи записи")
	urls := s.GetAllUser(context.Background(), userID)
	assert.Equal(t, map[string]string{id: url}, urls, "получение URL пользователя")
	urls = s.GetAllUser(context.Background(), userWithoutURLsID)
	assert.Equal(t, map[string]string{}, urls, "получение URL пользователя, не добавлявшего URL")
}
