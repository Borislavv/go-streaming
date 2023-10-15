package storage

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"github.com/Borislavv/video-streaming/internal/domain/logger"
	"github.com/Borislavv/video-streaming/internal/infrastructure/helper"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
)

type Filesystem struct {
	ctx    context.Context
	logger logger.Logger
}

func NewFilesystemStorage(ctx context.Context, logger logger.Logger) *Filesystem {
	return &Filesystem{
		ctx:    ctx,
		logger: logger,
	}
}

// Has is checking whether the file already exists.
func (s *Filesystem) Has(filename string) (has bool, e error) {
	// resources dir.
	resourcesDir, err := helper.ResourcesDir()
	if err != nil {
		return true, s.logger.LogPropagate(err)
	}

	// resources files dir.
	dir, err := os.Open(resourcesDir)
	if err != nil {
		return true, s.logger.LogPropagate(err)
	}
	defer func() { _ = dir.Close() }()

	// slice of string which is filenames
	filenames, err := dir.Readdirnames(-1)
	if err != nil {
		return true, s.logger.LogPropagate(err)
	}

	// attempt of finding a match
	for _, foundFilename := range filenames {
		if foundFilename == filename {
			return true, nil
		}
	}
	return false, nil
}

// Store is saving file and calculating new hashed name.
func (s *Filesystem) Store(
	name string,
	reader io.Reader,
) (
	length int64,
	filename string,
	filepath string,
	err error,
) {
	// resource file name
	filename = name

	// resources files directory
	dir, err := helper.ResourcesDir()
	if err != nil {
		return 0, "", "", s.logger.LogPropagate(err)
	}

	// full qualified file path
	filepath = fmt.Sprintf("%v%v", dir, name)

	// resource creating which will represented as a simple file at now
	createdFile, err := os.Create(filepath)
	if err != nil {
		return 0, "", "", s.logger.LogPropagate(err)
	}
	defer func() { _ = createdFile.Close() }()

	// moving the data in to the created file from tmp
	length, err = io.Copy(createdFile, reader)
	if err != nil {
		return 0, "", "", s.logger.LogPropagate(err)
	}

	// returning id of the created file, e.g. resourceId
	return length, filename, filepath, nil
}

func (s *Filesystem) Remove(name string) error {
	dir, err := helper.ResourcesDir()
	if err != nil {

	}
}

// getFilename - will return calculated filename with extension
func (s *Filesystem) getFilename(header *multipart.FileHeader) (filename string, e error) {
	hash := sha256.New()
	if _, err := hash.Write(
		[]byte(
			fmt.Sprintf(
				"%v%d%+v",
				header.Filename,
				header.Size,
				header.Header,
			),
		),
	); err != nil {
		return "", s.logger.LogPropagate(err)
	}

	return fmt.Sprintf(
		"%v%v",
		hex.EncodeToString(hash.Sum(nil)),
		filepath.Ext(header.Filename),
	), nil
}
