package file

import (
	"context"
	"fmt"
	"github.com/Borislavv/video-streaming/internal/domain/logger/interface"
	"github.com/Borislavv/video-streaming/internal/domain/service/di/interface"
	"github.com/Borislavv/video-streaming/internal/domain/vo"
	"github.com/Borislavv/video-streaming/internal/infrastructure/helper"
	"io"
	"os"
)

type FilesystemStorageService struct {
	ctx    context.Context
	logger loggerinterface.Logger
}

func NewFilesystemStorageService(serviceContainer diinterface.ServiceContainer) (*FilesystemStorageService, error) {
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
func (s *FilesystemStorageService) Has(userID vo.ID, filename string) (has bool, err error) {
	// resources dir.
	resourcesDir, err := helper.ResourcesDir()
	if err != nil {
		return true, s.logger.LogPropagate(err)
	}

	// dir. with other dirs. which separates by userIDs
	usersDirs, err := os.Open(resourcesDir)
	if err != nil {
		return true, s.logger.LogPropagate(err)
	}
	defer func() { _ = usersDirs.Close() }()

	// dirs. by userIDs (name of each dir is userID)
	userDirs, err := usersDirs.Readdirnames(-1)
	if err != nil {
		return true, s.logger.LogPropagate(err)
	}

	// attempt of finding a match
	for _, userIDHex := range userDirs {
		if userIDHex == userID.Hex() {
			userDir, err := os.Open(resourcesDir + userIDHex)
			if err != nil {
				return true, s.logger.LogPropagate(err)
			}

			userFiles, err := userDir.Readdirnames(-1)
			if err != nil {
				return true, s.logger.LogPropagate(err)
			}

			for _, userFilename := range userFiles {
				if userFilename == filename {
					return true, nil
				}
			}
		}
	}

	return false, nil
}

// Store is saving file and calculating new hashed name.
func (s *FilesystemStorageService) Store(
	userID vo.ID,
	filename string,
	reader io.Reader,
) (
	length int64,
	path string,
	err error,
) {
	// resources files directory
	dir, err := helper.ResourcesDir()
	if err != nil {
		return 0, "", s.logger.LogPropagate(err)
	}

	// user files directory path
	dir = fmt.Sprintf("%v%v", dir, userID.Hex())

	// create a user directory if not exists
	if _, err = os.Stat(dir); os.IsNotExist(err) {
		if err = os.Mkdir(dir, 0644); err != nil {
			return 0, "", s.logger.LogPropagate(err)
		}
	}

	// full qualified file path
	path = fmt.Sprintf("%v/%v", dir, filename)

	// resource creating which will represented as a simple file at now
	createdFile, err := os.Create(path)
	if err != nil {
		return 0, "", s.logger.LogPropagate(err)
	}
	defer func() { _ = createdFile.Close() }()

	// moving the data in to the created file from tmp
	length, err = io.Copy(createdFile, reader)
	if err != nil {
		return 0, "", s.logger.LogPropagate(err)
	}

	// returning id of the created file, e.g. resourceId
	return length, path, nil
}

func (s *FilesystemStorageService) Remove(userID vo.ID, name string) error {
	// full qualified filepath
	path, err := s.filepath(userID, name)
	if err != nil {
		return s.logger.LogPropagate(err)
	}

	// file already is not exists
	if _, err = os.Stat(path); os.IsNotExist(err) {
		return nil
	}

	// removing the target file
	if err = os.Remove(path); err != nil {
		return s.logger.LogPropagate(err)
	}

	return nil
}

func (s *FilesystemStorageService) filepath(userID vo.ID, filename string) (path string, err error) {
	// resources dir path
	dir, err := helper.ResourcesDir()
	if err != nil {
		return "", s.logger.LogPropagate(err)
	}

	// user files dir path
	dir = fmt.Sprintf("%v%v", dir, userID.Hex())

	// create a user directory if not exists
	if _, err = os.Stat(dir); os.IsNotExist(err) {
		if err = os.Mkdir(dir, 0644); err != nil {
			return "", s.logger.LogPropagate(err)
		}
	}

	return fmt.Sprintf("%v/%v", dir, filename), nil
}
