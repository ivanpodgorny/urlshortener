package storage

import (
	"bufio"
	"context"
	"os"
	"strings"
)

const (
	urlSectionName  = "url"
	userSectionName = "user"
)

// Memory реализует интерфейс service.Storage для хранения url в памяти.
// Если передать в конструктор файловый дескриптор, будет также сохранять
// url в открытый файл.
type Memory struct {
	urls       map[string]string
	userData   map[string][]string
	persistent *os.File
}

func NewMemory(file *os.File) *Memory {
	s := Memory{
		urls:       map[string]string{},
		userData:   map[string][]string{},
		persistent: file,
	}
	s.loadDataInMemory()

	return &s
}

func (f *Memory) Add(_ context.Context, id string, url string, userID string) (string, error) {
	if _, exist := f.urls[id]; exist {
		return "", ErrKeyExists
	}

	if err := f.saveToPersistent(urlSectionName, id, url); err != nil {
		return "", err
	}
	if err := f.saveToPersistent(userSectionName, userID, id); err != nil {
		return "", err
	}
	f.urls[id] = url
	f.userData[userID] = append(f.userData[userID], id)

	return id, nil
}

func (f *Memory) Get(_ context.Context, id string) (string, error) {
	url, ok := f.urls[id]
	if !ok {
		return "", ErrKeyNotFound
	}

	return url, nil
}

func (f *Memory) GetAllUser(ctx context.Context, userID string) map[string]string {
	data := map[string]string{}
	ids, ok := f.userData[userID]
	if !ok {
		return data
	}

	for _, id := range ids {
		url, err := f.Get(ctx, id)
		if err == nil {
			data[id] = url
		}
	}

	return data
}

func (f *Memory) loadDataInMemory() {
	if f.persistent == nil {
		return
	}

	scanner := bufio.NewScanner(f.persistent)
	for scanner.Scan() {
		sectionAndKeyVal := strings.Split(scanner.Text(), ",")
		switch sectionAndKeyVal[0] {
		case urlSectionName:
			f.urls[sectionAndKeyVal[1]] = sectionAndKeyVal[2]
		case userSectionName:
			f.userData[sectionAndKeyVal[1]] = append(f.userData[sectionAndKeyVal[1]], sectionAndKeyVal[2])
		}
	}
}

func (f *Memory) saveToPersistent(section, key, val string) error {
	if f.persistent == nil {
		return nil
	}

	_, err := f.persistent.Write([]byte(section + "," + key + "," + val + "\n"))

	return err
}
