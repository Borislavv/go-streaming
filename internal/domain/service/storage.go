package service

import (
	"mime/multipart"
)

type Storage interface {
	// Has is checking whether the file already exists.
	Has(header *multipart.FileHeader) (has bool, e error)
	Store(file multipart.File, header *multipart.FileHeader) (id string, e error)
}
