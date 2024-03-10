package file

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"github.com/Borislavv/video-streaming/internal/domain/logger/interface"
	"github.com/Borislavv/video-streaming/internal/domain/service/di/interface"
	"github.com/Borislavv/video-streaming/internal/infrastructure/helper"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
)

type FilesystemStorageService struct {
	ctx    context.Context
	logger loggerinterface.Logger
}

func NewFilesystemStorageService(serviceContainer diinterface.ContainerManager) (*FilesystemStorageService, error) {
	loggerService, err := serviceContainer.GetLoggerService()
	if err != nil {
		return nil, err
	}

	ctx, err := serviceContainer.GetCtx()
	if err != nil {
		return nil, loggerService.LogPropagate(err)
	}

	return &FilesystemStorageService{
		ctx:    ctx,
		logger: loggerService,
	}, nil
}

// Has is checking whether the file already exists.
func (s *FilesystemStorageService) Has(filename string) (has bool, e error) {
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
func (s *FilesystemStorageService) Store(
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

func (s *FilesystemStorageService) Remove(name string) error {
	// resources files directory
	dir, err := helper.ResourcesDir()
	if err != nil {
		return s.logger.LogPropagate(err)
	}

	// full qualified file path
	path := fmt.Sprintf("%v%v", dir, name)

	// removing the target file
	if err = os.Remove(path); err != nil {
		return s.logger.LogPropagate(err)
	}

	return nil
}

// getFilename - will return calculated filename with extension
func (s *FilesystemStorageService) getFilename(header *multipart.FileHeader) (filename string, err error) {
	hash := sha256.New()
	if _, err = hash.Write(
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
