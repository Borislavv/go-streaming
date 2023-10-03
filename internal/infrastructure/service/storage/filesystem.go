package storage

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"github.com/Borislavv/video-streaming/internal/domain/logger"
	"github.com/Borislavv/video-streaming/internal/infrastructure/helper"
	"io"
	"log"
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

	// moving the data in to the created file from tmp
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

	wg := &sync.WaitGroup{}
	ctx, cancel := context.WithCancel(s.ctx)
	dataCh := make(chan []byte)

	chunkSize := 1024 * 1024 * 1

	taskProvidersNum := 1
	dataProvidersNum := runtime.NumCPU()
	dataConsumersNum := 1

	taskCh := make(chan struct{}, dataProvidersNum)

	// taskProvider: sending task to each provider for control the treads state
	wg.Add(taskProvidersNum)
	go func() {
		defer wg.Done()
		for {
			select {
			case <-ctx.Done():
				close(taskCh)
				log.Printf("closed taskCh, interruted by error")
				return
			default:
				taskCh <- struct{}{}
			}
		}
	}()

	wg2 := &sync.WaitGroup{}
	once := &sync.Once{}

	// dataProviders: consuming tasks, reading the file and send the slices to data consumer
	wg.Add(dataProvidersNum)
	go func() {
		for i := 0; i < dataProvidersNum; i++ {
			go func() {
				defer wg.Done()

				for range taskCh {
					buff := make([]byte, 0, 1024*1024*1)
					n, err := part.Read(buff)
					if err != nil {
						if err == io.EOF {
							return
						}
						s.logger.Critical(err)
						return
					}
					if n < chunkSize { // handle the last chunk
						if n == 0 {
							log.Printf("dataProvider: readed 0 bytes, RETURNED")

							once.Do(
								func() {
									wg2.Add(1)
									defer wg2.Done()
									cancel()
									close(dataCh)
									wg.Wait()
								},
							)

							return
						}
						dataCh <- buff[:n]
						return
					}
					log.Printf("dataProvider: readed %d bytes", n)
					dataCh <- buff
				}
			}()
		}
	}()

	wg2.Add(dataConsumersNum)
	go func() {
		defer wg2.Done()

		for data := range dataCh {
			n, err := createdFile.Write(data)
			if err != nil {
				s.logger.Critical(err)
				log.Println(err)

				cancel() // close the tasks provider

				wg2.Add(1)
				go func() { // skipping the last data if err occurred
					defer wg2.Done()
					for range dataCh {
					} // deadlock escaping
				}()
				wg.Wait() // wait while the previous goroutines will be closed
				close(dataCh)
			}
			length += int64(n)
		}
	}()
	wg2.Wait()

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
