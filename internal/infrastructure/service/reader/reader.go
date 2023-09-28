package reader

import (
	"context"
	"fmt"
	"github.com/Borislavv/video-streaming/internal/domain/logger"
	"github.com/Borislavv/video-streaming/internal/infrastructure/service/reader/model"
	"os"
)

const chunksChBuffer = 1

type FileReaderService struct {
	ctx       context.Context
	logger    logger.Logger
	chunkSize int
}

func NewFileReaderService(ctx context.Context, logger logger.Logger, chunkSize int) *FileReaderService {
	return &FileReaderService{ctx: ctx, logger: logger, chunkSize: chunkSize}
}

func (r *FileReaderService) ReadAll(file *os.File) *model.Chunk {
	return model.NewChunk(0)
}

// ReadByChunks - reads a file by separated chunks
// and passed it into the channel (chunk size is setting up through env. configuration).
func (r *FileReaderService) ReadByChunks(file *os.File, offset int64) chan *model.Chunk {
	r.logger.Info(fmt.Sprintf("file '%v' reading started", file.Name()))

	stat, err := file.Stat()
	if err != nil {
		r.logger.Info(fmt.Sprintf("file '%v' reading file stat with errors: %v", file.Name(), err))
		return nil
	}

	ch := make(chan *model.Chunk, chunksChBuffer)
	go func(offset int64, size int64) {
		defer close(ch)
		for {
			select {
			case <-r.ctx.Done():
				r.logger.Info(fmt.Sprintf("file '%v' reading interrupted", file.Name()))
				return
			default:
				// compute the current chunk buffer
				currentChunkSize := r.chunkSize
				if int64(currentChunkSize) > (size - offset) {
					currentChunkSize = int(size - offset)
					if currentChunkSize == 0 {
						r.logger.Info(fmt.Sprintf("file '%v' reading finished properly", file.Name()))
						return
					}
				}

				// make a new chunk with appropriate buffer
				chunk := model.NewChunk(currentChunkSize)

				// read the current batch of bites
				chunk.Len, err = file.ReadAt(chunk.Data, offset)
				if err != nil {
					r.logger.Error(err)
					r.logger.Info(fmt.Sprintf("file '%v' reading finished with errors", file.Name()))
					return
				}
				offset += int64(chunk.Len)

				// cut the last chunk to its real length
				if chunk.GetLen() < r.chunkSize {
					chunk.SetData(chunk.GetData()[:chunk.GetLen()])
				}

				// sent the chunk to consumer
				ch <- chunk
			}
		}
	}(offset, stat.Size())
	return ch
}
