package storage

import (
	"io"
)

type Storage interface {
	// Has is checking whether the file already exists.
	Has(filename string) (has bool, err error)
	// Store is saving file and calculating new hashed name.
	Store(name string, reader io.Reader) (length int64, filename string, filepath string, err error)
	// Remove is delete the file by name from resources directory.
	Remove(name string) (err error)
}
