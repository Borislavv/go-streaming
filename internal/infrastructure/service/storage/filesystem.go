package storage

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"github.com/Borislavv/video-streaming/internal/domain/logger"
	"github.com/Borislavv/video-streaming/internal/infrastructure/helper"
	"io"
	"math"
	"mime/multipart"
	"os"
	"path/filepath"
	"sync"
)

type Filesystem struct {
	logger logger.Logger
}

func NewFilesystemStorage(logger logger.Logger) *Filesystem {
	return &Filesystem{
		logger: logger,
	}
}

// Has is checking whether the file already exists.
func (s *Filesystem) Has(header *multipart.FileHeader) (has bool, e error) {
	// filename with extension
	name, err := s.getFilename(header)
	if err != nil {
		return true, err
	}

	resourcesDir, err := helper.ResourcesDir()
	if err != nil {
		return true, err
	}

	// resources files directory
	dir, err := os.Open(resourcesDir)
	if err != nil {
		return true, err
	}
	defer func() { _ = dir.Close() }()

	// slice of string which is filenames
	filenames, err := dir.Readdirnames(-1)
	if err != nil {
		return true, err
	}

	// attempt of finding a match
	for _, filename := range filenames {
		if filename == name {
			return true, nil
		}
	}
	return false, nil
}

func (s *Filesystem) Store(file multipart.File, header *multipart.FileHeader) (name string, path string, err error) {
	// filename with extension
	name, err = s.getFilename(header)
	if err != nil {
		return "", "", err
	}

	// resources files directory
	dir, err := helper.ResourcesDir()
	if err != nil {
		return "", "", err
	}

	// full qualified file path
	path = fmt.Sprintf("%v%v", dir, name)

	// resource creating which will represented as a simple file at now
	createdFile, err := os.Create(path)
	if err != nil {
		return "", "", err
	}
	defer func() { _ = createdFile.Close() }()

	// moving the data in to the created file from tmp
	_, err = io.Copy(createdFile, file)
	if err != nil {
		return "", "", err
	}

	// returning id of the created file, e.g. resourceId
	return name, path, nil
}

func (s *Filesystem) StoreConcurrently(file multipart.File, header *multipart.FileHeader) (name string, path string, err error) {
	// filename with extension
	name, err = s.getFilename(header)
	if err != nil {
		return "", "", err
	}

	// resources files directory
	dir, err := helper.ResourcesDir()
	if err != nil {
		return "", "", err
	}

	// full qualified file path
	path = fmt.Sprintf("%v%v", dir, name)

	// resource creating which will represented as a simple file at now
	createdFile, err := os.Create(path)
	if err != nil {
		return "", "", err
	}
	defer func() { _ = createdFile.Close() }()

	threads := 5
	fileSize := header.Size
	chunkSize := int64(1024 * 1024 * 1) // 1mb.
	chunksNumber := int64(math.Ceil(float64(fileSize / chunkSize)))

	wg := &sync.WaitGroup{}
	taskCh := make(chan *struct {
		num    int64
		offset int64
	}, threads)

	wg.Add(1)
	go func() {
		defer wg.Done()

		for i := int64(0); i < chunksNumber; i++ {
			taskCh <- &struct {
				num    int64
				offset int64
			}{num: i, offset: i * chunkSize}
		}
	}()

	n, err := file.ReadAt([]byte{}, 0)
	if err != nil {
		return "", "", err
	}

	createdFile.WriteAt()
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
		return "", err
	}

	return fmt.Sprintf(
		"%v%v",
		hex.EncodeToString(hash.Sum(nil)),
		filepath.Ext(header.Filename),
	), nil
}
