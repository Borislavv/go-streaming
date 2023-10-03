package storage

import (
	"mime/multipart"
)

type Storage interface {
	// Has is checking whether the file already exists.
	Has(filename string) (has bool, err error)
	// Store is saving file and calculating new hashed name.
	Store(name string, part *multipart.Part) (length int64, filename string, filepath string, err error)
	// StoreConcurrently is saving file and calculating new hashed name concurrently.
	StoreConcurrently(name string, part *multipart.Part) (length int64, filename string, filepath string, err error)
}
