package storage

import (
	"bufio"
	"context"
	"errors"
	"io"
	"os"
	"strings"
	"sync"

	inerr "github.com/ivanpodgorny/urlshortener/internal/app/errors"
)

// Memory реализует интерфейс service.Storage для хранения url в памяти.
// Если передать в конструктор файловый дескриптор, будет также сохранять
// url в открытый файл.
type Memory struct {
	urls       map[string]string
	userData   map[string][]string
	persistent *os.File
	mu         sync.RWMutex
}

const (
	deletedFlag     = "deleted"
	urlSectionName  = "url"
	userSectionName = "user"
)

// ErrKeyExists URL с данным id eже существует.
var ErrKeyExists = errors.New("key exists")

// ErrKeyNotFound н найден URL с данным id.
var ErrKeyNotFound = errors.New("key not found")

// NewMemory возвращает указатель на новый экземпляр Memory.
func NewMemory(file *os.File) *Memory {
	s := Memory{
		urls:       map[string]string{},
		userData:   map[string][]string{},
		persistent: file,
	}
	s.loadDataInMemory()

	return &s
}

// Add сохраняет URL.
func (m *Memory) Add(_ context.Context, id string, url string, userID string) (string, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, exist := m.urls[id]; exist {
		return "", ErrKeyExists
	}

	if err := m.saveToPersistent(urlSectionName, id, url); err != nil {
		return "", err
	}
	if err := m.saveToPersistent(userSectionName, userID, id); err != nil {
		return "", err
	}
	m.urls[id] = url
	m.userData[userID] = append(m.userData[userID], id)

	return id, nil
}

// Get возвращает сохраненный URL по id. Если URL был помечен удаленным, возвращает
// ошибку errors.ErrURLIsDeleted.
func (m *Memory) Get(_ context.Context, id string) (string, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	url, ok := m.urls[id]
	if !ok {
		return "", ErrKeyNotFound
	}

	if url == deletedFlag {
		return "", inerr.ErrURLIsDeleted
	}

	return url, nil
}

// GetAllUser возвращает все сохраненные URL пользователя.
func (m *Memory) GetAllUser(ctx context.Context, userID string) map[string]string {
	m.mu.RLock()
	defer m.mu.RUnlock()

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

// DeleteBatch удаляет URL с заданными id.
func (m *Memory) DeleteBatch(_ context.Context, urlIDs []string, userID string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	for _, urlID := range urlIDs {
		if !m.belongsToUser(urlID, userID) {
			continue
		}

		m.urls[urlID] = deletedFlag
	}

	return m.renewPersistent()
}

func (m *Memory) loadDataInMemory() {
	if m.persistent == nil {
		return
	}

	scanner := bufio.NewScanner(m.persistent)
	for scanner.Scan() {
		sectionAndKeyVal := strings.Split(scanner.Text(), ",")
		switch sectionAndKeyVal[0] {
		case urlSectionName:
			m.urls[sectionAndKeyVal[1]] = sectionAndKeyVal[2]
		case userSectionName:
			m.userData[sectionAndKeyVal[1]] = append(m.userData[sectionAndKeyVal[1]], sectionAndKeyVal[2])
		}
	}
}

func (m *Memory) renewPersistent() error {
	if m.persistent == nil {
		return nil
	}

	if _, err := m.persistent.Seek(0, io.SeekStart); err != nil {
		return err
	}
	if err := m.persistent.Truncate(0); err != nil {
		return err
	}

	for id, url := range m.urls {
		if err := m.saveToPersistent(urlSectionName, id, url); err != nil {
			return err
		}
	}
	for userID, ids := range m.userData {
		for _, id := range ids {
			if err := m.saveToPersistent(userSectionName, userID, id); err != nil {
				return err
			}
		}
	}

	return nil
}

func (m *Memory) saveToPersistent(section, key, val string) error {
	if m.persistent == nil {
		return nil
	}

	_, err := m.persistent.Write([]byte(section + "," + key + "," + val + "\n"))

	return err
}

func (m *Memory) belongsToUser(urlID, userID string) bool {
	for _, id := range m.userData[userID] {
		if id == urlID {
			return true
		}
	}

	return false
}
