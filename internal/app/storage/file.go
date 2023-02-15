package storage

import (
	"bufio"
	"context"
	"os"
	"strings"
)

// File реализует интерфейс service.Storage
// для хранения url в файле.
type File struct {
	memory     map[string]string
	persistent *os.File
}

func NewFile(filename string) (*File, error) {
	file, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0600)
	if err != nil {
		return nil, err
	}

	s := File{
		memory:     map[string]string{},
		persistent: file,
	}
	s.loadDataInMemory()

	return &s, nil
}

func (f *File) Add(_ context.Context, id string, url string) error {
	if _, exist := f.memory[id]; exist {
		return ErrKeyExists
	}

	if err := f.saveToPersistent(id, url); err != nil {
		return err
	}
	f.memory[id] = url

	return nil
}

func (f *File) Get(_ context.Context, id string) (string, error) {
	url, ok := f.memory[id]
	if !ok {
		return "", ErrKeyNotFound
	}

	return url, nil
}

func (f *File) Close() error {
	return f.persistent.Close()
}

func (f *File) loadDataInMemory() {
	scanner := bufio.NewScanner(f.persistent)
	for scanner.Scan() {
		idAndURL := strings.Split(scanner.Text(), ",")
		f.memory[idAndURL[0]] = idAndURL[1]
	}
}

func (f *File) saveToPersistent(id, url string) error {
	_, err := f.persistent.Write([]byte(id + "," + url + "\n"))

	return err
}
