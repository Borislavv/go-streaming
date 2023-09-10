package filesystem

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"github.com/Borislavv/video-streaming/internal/infrastructure/helper"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
)

type Storage struct {
}

func NewStorage() *Storage {
	return &Storage{}
}

// Has is checking whether the file already exists.
func (s *Storage) Has(header *multipart.FileHeader) (has bool, e error) {
	resourceId, err := s.hashSum(header)
	if err != nil {
		return true, err
	}

	resourcesDir, err := helper.ResourcesDir()
	if err != nil {
		return true, err
	}

	dir, err := os.Open(resourcesDir)
	if err != nil {
		return true, err
	}
	defer dir.Close()

	fileNames, err := dir.Readdirnames(-1)
	if err != nil {
		return true, err
	}

	for _, fileName := range fileNames {
		if strings.HasPrefix(fileName, resourceId) {
			return true, nil
		}
	}
	return false, nil
}

func (s *Storage) Store(file multipart.File, header *multipart.FileHeader) (id string, e error) {
	// filename without extension
	id, err := s.hashSum(header)
	if err != nil {
		return "", err
	}

	// filename with extension
	name := fmt.Sprintf("%v%v", id, filepath.Ext(header.Filename))

	// resources files directory
	dir, err := helper.ResourcesDir()
	if err != nil {
		return "", err
	}

	// full qualified file path
	path := fmt.Sprintf("%v%v", dir, name)

	// resource creating which will represented as a simple file at now
	createdFile, err := os.Create(path)
	if err != nil {
		return "", err
	}
	defer createdFile.Close()

	// moving the data in to the created file from tmp
	_, err = io.Copy(createdFile, file)
	if err != nil {
		return "", err
	}

	// returning id of the created file
	return id, nil
}

func (s *Storage) hashSum(header *multipart.FileHeader) (id string, e error) {
	hash := sha256.New()
	if _, err := hash.Write(
		[]byte(
			fmt.Sprintf("%v%d%+v", header.Filename, header.Size, header.Header),
		),
	); err != nil {
		return "", err
	}
	return hex.EncodeToString(hash.Sum(nil)), nil
}
