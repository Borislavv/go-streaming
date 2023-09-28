package reader

import (
	"context"
	"fmt"
	"github.com/Borislavv/video-streaming/internal/domain/logger"
	"github.com/Borislavv/video-streaming/internal/infrastructure/service/reader/model"
	"os"
	"sync"
)

const (
	chunksChBuffer = 1
	readingThreads = 5
)

type FileReaderService struct {
	ctx       context.Context
	logger    logger.Logger
	chunkSize int
}

func NewFileReaderService(ctx context.Context, logger logger.Logger, chunkSize int) *FileReaderService {
	return &FileReaderService{ctx: ctx, logger: logger, chunkSize: chunkSize}
}

// TODO Must be tested!
// ReadAll - reads a whole file in a single chunk.
func (r *FileReaderService) ReadAll(file *os.File) *model.Chunk {
	r.logger.Info(fmt.Sprintf("reading all file '%v' started", file.Name()))

	stat, err := file.Stat()
	if err != nil {
		r.logger.Critical(fmt.Sprintf("reading all file '%v' error: %v", file.Name(), err))
		return nil
	}

	// chunks number
	chunks := stat.Size() / int64(r.chunkSize)
	// reading threads number
	threads := int64(readingThreads)
	// check the num of chunks more than threads
	if threads > chunks {
		threads = chunks
	}

	wg := &sync.WaitGroup{}

	taskCh := make(chan *struct {
		num    int64
		offset int64
		length int64
	}, threads)

	// provider
	wg.Add(1)
	go func() {
		defer wg.Done()
		defer close(taskCh)
		// make a task for each chunk and send to consumers
		for chk := int64(0); chk < chunks; chk++ {
			offset := chk * int64(r.chunkSize)

			length := int64(r.chunkSize)
			if length > (stat.Size() - offset) {
				length = stat.Size() - offset
			}

			taskCh <- &struct {
				num    int64
				offset int64
				length int64
			}{
				num:    chk,
				offset: offset,
				length: length,
			}
		}
	}()

	// map of file chunks by int64 serial number
	fileMap := make(map[int64][]byte, chunks)
	// make a mutex for concurrently write into the file map
	mu := &sync.Mutex{}

	// consumer
	wg.Add(int(threads))
	go func() {
		for thrd := int64(0); thrd < threads; thrd++ {
			go func(thrd int64) {
				defer wg.Done()

				for task := range taskCh {
					buff := make([]byte, task.length)
					_, err := file.ReadAt(buff, task.offset)
					if err != nil {
						r.logger.Critical(
							fmt.Sprintf("reading all file '%v' error: %v at %d thread", file.Name(), err, thrd),
						)
						return
					}

					mu.Lock()
					fileMap[task.num] = buff
					mu.Unlock()
				}
			}(thrd)
		}
	}()

	// awaiting while whole file will be read
	wg.Wait()

	// collect the entire file into one chunk
	chunk := model.NewChunk(stat.Size())
	for i := int64(0); i < int64(len(fileMap)); i++ {
		chunk.Data = append(chunk.Data, fileMap[i]...)
	}
	return chunk
}

// ReadByChunks - reads a file by separated chunks
// and passed it into the channel (chunk size is setting up through env. configuration).
func (r *FileReaderService) ReadByChunks(file *os.File, offset int64) chan *model.Chunk {
	r.logger.Info(fmt.Sprintf("reading file '%v' by chunks started", file.Name()))

	stat, err := file.Stat()
	if err != nil {
		r.logger.Info(fmt.Sprintf("reading file '%v' by chunks file stat with errors: %v", file.Name(), err))
		return nil
	}

	ch := make(chan *model.Chunk, chunksChBuffer)
	go func(offset int64, size int64) {
		defer close(ch)
		for {
			select {
			case <-r.ctx.Done():
				r.logger.Info(fmt.Sprintf("reading file '%v' by chunks interrupted", file.Name()))
				return
			default:
				// compute the current chunk buffer
				currentChunkSize := int64(r.chunkSize)
				if currentChunkSize > (size - offset) {
					currentChunkSize = size - offset
					if currentChunkSize == 0 {
						r.logger.Info(fmt.Sprintf("reading file '%v' by chunks finished properly", file.Name()))
						return
					}
				}

				// make a new chunk with appropriate buffer
				chunk := model.NewChunk(currentChunkSize)

				// read the current batch of bites
				length, err := file.ReadAt(chunk.Data, offset)
				if err != nil {
					r.logger.Error(err)
					r.logger.Info(fmt.Sprintf("reading file '%v' by chunks finished with errors", file.Name()))
					return
				}
				chunk.SetLen(int64(length))
				offset += chunk.GetLen()

				// cut the last chunk to its real length
				if chunk.GetLen() < int64(r.chunkSize) {
					chunk.SetData(chunk.GetData()[:chunk.GetLen()])
				}

				// sent the chunk to consumer
				ch <- chunk
			}
		}
	}(offset, stat.Size())
	return ch
}
