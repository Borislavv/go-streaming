package storage

import (
	"mime/multipart"
)

type Storage interface {
	// Has is checking whether the file already exists.
	Has(filename string) (has bool, err error)

	// TODO This both methods must be merged
	// Store is saving file and calculating new hashed name.
	Store(name string, part *multipart.Part) (length int64, filename string, filepath string, err error)
	// StoreFile is saving file and calculating new hashed name.
	StoreFile(name string, part *multipart.File) (length int64, filename string, filepath string, err error)
}
