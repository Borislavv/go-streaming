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
	"runtime"
	"sync"
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

	//// moving the data in to the created file from tmp
	length, err = io.Copy(createdFile, part)
	if err != nil {
		return 0, "", "", err
	}

	// returning id of the created file, e.g. resourceId
	return length, filename, filepath, nil
}

func (s *Filesystem) StoreConcurrently(
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

	var wg sync.WaitGroup
	var wgp sync.WaitGroup

	dataCh := make(chan []byte)
	chunkSize := 1024 * 1024 * 1
	dataProvidersNum := runtime.NumCPU()
	doneCh := make(chan struct{})

	wg.Add(dataProvidersNum)
	wgp.Add(dataProvidersNum)
	for i := 0; i < dataProvidersNum; i++ {
		go func() {
			defer func() {
				wg.Done()
				wgp.Done()
			}()

			for {
				select {
				case <-doneCh:
					s.logger.Critical("reading interrupted")
					return
				default:
					buff := make([]byte, 1024*1024*50) // TODO try to reduce to 4096 if the file will be built successfully
					n, e := part.Read(buff)
					if e != nil && e != io.EOF {
						s.logger.Critical(e)
						return
					}
					s.logger.Critical(fmt.Sprintf("found %d bytes and ent through dataCh", n))
					if n < 1024*1024*50 {
						if n == 0 {
							s.logger.Critical("zero bytes found, exit")
							return // normal exit
						}
						dataCh <- buff[:n]
						s.logger.Critical("found slice of bytes which is lower than chunkSize")
						return // normal exit
					}
					s.logger.Emergency("buffer is MORE THAN 5000 bytes!!!!!!11111111111111")
					dataCh <- buff
				}
			}
		}()
	}

	wg.Add(1)
	go func() {
		defer wg.Done()
		wgp.Wait()
		close(dataCh)
		s.logger.Critical("dataCh is closed")
	}()

	wg.Add(1)
	go func() {
		defer func() {
			wg.Done()
			s.logger.Critical("main consumer is closed")
		}()

		buff := make([]byte, chunkSize)
		buffLen := 0

		for data := range dataCh {
			if buffLen+len(data) > chunkSize {
				// flush the buffer
				n, e := createdFile.Write(buff[:buffLen])
				if e != nil {
					s.logger.Critical(e)
					close(doneCh)
					wg.Add(1)
					go func() {
						defer func() {
							wg.Done()
							s.logger.Critical("child consumer is closed")
						}()
						for range dataCh {
						}
					}()
					err = e
					return
				}

				buff = buff[:0]
				buff = append(buff, data...)
				buffLen = len(data)

				s.logger.Info(fmt.Sprintf("wrote %d bytes", n))
				length += int64(n)
			} else {
				buff = append(buff, data...)
				buffLen += len(data)
			}
		}

		if buffLen > 0 {
			n, e := createdFile.Write(buff[:buffLen])
			if e != nil {
				s.logger.Critical(e)
				close(doneCh)
				wg.Add(1)
				go func() {
					defer func() {
						wg.Done()
						s.logger.Critical("child consumer is closed")
					}()
					for range dataCh {
					}
				}()
				err = e
				return
			}
			s.logger.Info(fmt.Sprintf("wrote %d bytes", n))
			length += int64(n)
		}
	}()

	wg.Wait()

	s.logger.Info(fmt.Sprintf("%d %v %v %v", length, filename, filepath, err))

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
