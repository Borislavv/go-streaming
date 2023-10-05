package storage

import (
	"bufio"
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
		return true, err
	}

	// resources files dir.
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
	part *multipart.Part,
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
		return 0, "", "", err
	}

	// full qualified file path
	filepath = fmt.Sprintf("%v%v", dir, name)

	// resource creating which will represented as a simple file at now
	createdFile, err := os.Create(filepath)
	if err != nil {
		return 0, "", "", err
	}
	defer func() { _ = createdFile.Close() }()

	// moving the data in to the created file from tmp
	length, err = io.Copy(createdFile, part)
	if err != nil {
		return 0, "", "", err
	}

	// returning id of the created file, e.g. resourceId
	return length, filename, filepath, nil
}

// StoreBuffered is saving file and calculating new hashed name (buffered).
func (s *Filesystem) StoreBuffered(
	name string,
	part *multipart.Part,
) (
	length int64,
	filename string,
	filepath string,
	err error,
) {
	filename = name

	// resources files directory
	dir, err := helper.ResourcesDir()
	if err != nil {
		return 0, "", "", err
	}

	// full qualified file path
	filepath = fmt.Sprintf("%v%v", dir, name)

	// resource creating which will represented as a simple file at now
	createdFile, err := os.Create(filepath)
	if err != nil {
		return 0, "", "", err
	}
	defer func() { _ = createdFile.Close() }()

	buf := bufio.NewReader(part)
	reader := io.MultiReader(buf, io.LimitReader(part, 1024*1024*1024*10))

	chunkLen := 0
	chunkBuff := make([]byte, 1024*1024*3)
	for {
		buff := make([]byte, 1024*1024)
		r, e := reader.Read(buff)
		if e != nil && e != io.EOF {
			s.logger.Critical(e)
			err = e
			break
		}
		s.logger.Debug(fmt.Sprintf("Read %d bytes from multipart.Part", r))
		if chunkLen+r > len(chunkBuff) || r == 0 {
			// flush chunk buffer
			w, e := createdFile.Write(chunkBuff[:chunkLen])
			if e != nil {
				s.logger.Critical(e)
				err = e
				break
			}
			length += int64(w)

			chunkLen = 0
			chunkBuff = chunkBuff[:0]

			if r == 0 {
				break
			}
		} else {
			// append iteration buffer
			chunkBuff = append(chunkBuff, buff[:r]...)
			chunkLen += r
		}
	}

	if err != nil {
		return 0, "", "", err
	}

	return length, filename, filepath, nil
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
