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

// File реализует интерфейс service.Storage
// для хранения url в файле.
type File struct {
	urls       map[string]string
	userData   map[string][]string
	persistent *os.File
}

func NewFile(filename string) (*File, error) {
	file, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0600)
	if err != nil {
		return nil, err
	}

	s := File{
		urls:       map[string]string{},
		userData:   map[string][]string{},
		persistent: file,
	}
	s.loadDataInMemory()

	return &s, nil
}

func (f *File) Add(_ context.Context, id string, url string, userID string) error {
	if _, exist := f.urls[id]; exist {
		return ErrKeyExists
	}

	if err := f.saveToPersistent(urlSectionName, id, url); err != nil {
		return err
	}
	if err := f.saveToPersistent(userSectionName, userID, id); err != nil {
		return err
	}
	f.urls[id] = url
	f.userData[userID] = append(f.userData[userID], id)

	return nil
}

func (f *File) Get(_ context.Context, id string) (string, error) {
	url, ok := f.urls[id]
	if !ok {
		return "", ErrKeyNotFound
	}

	return url, nil
}

func (f *File) GetAllUser(ctx context.Context, userID string) map[string]string {
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

func (f *File) Close() error {
	return f.persistent.Close()
}

func (f *File) loadDataInMemory() {
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

func (f *File) saveToPersistent(section, key, val string) error {
	_, err := f.persistent.Write([]byte(section + "," + key + "," + val + "\n"))

	return err
}
