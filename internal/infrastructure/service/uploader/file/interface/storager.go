package fileinterface

import (
	"github.com/Borislavv/video-streaming/internal/domain/vo"
	"io"
)

type Storage interface {
	// Has is checking whether the file already exists.
	Has(userID vo.ID, filename string) (has bool, err error)
	// Store is saving file and calculating new hashed name.
	Store(userID vo.ID, name string, reader io.Reader) (length int64, filepath string, err error)
	// Remove is delete the file by name from resources directory.
	Remove(userID vo.ID, name string) (err error)
}
