package storage

import (
	"io"
)

type Storage interface {
	List() ([]string, error)
	FileExists(string) (bool, error)
	Get(string) (string, error)
	Put(io.Reader, string) error
	PutIfNotExists(io.Reader, string) error
	Delete(string) error
	DeleteList([]string) error
}
