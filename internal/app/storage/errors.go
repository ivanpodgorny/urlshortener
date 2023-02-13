package storage

import "errors"

var (
	ErrKeyExists   = errors.New("key exists")
	ErrKeyNotFound = errors.New("key not found")
)
