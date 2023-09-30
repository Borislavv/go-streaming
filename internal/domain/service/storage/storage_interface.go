package storage

import (
	"mime/multipart"
)

type Storage interface {
	// Has is checking whether the file already exists.
	Has(header *multipart.FileHeader) (has bool, err error)
	// Store is saving file and calculating new hashed name.
	Store(file multipart.File, header *multipart.FileHeader) (filename string, filepath string, err error)
}
