package read

import (
	"fmt"
	"github.com/Borislavv/video-streaming/internal/app/service/logger"
	"github.com/Borislavv/video-streaming/internal/domain/entity"
	"io"
	"os"
)

// ChunkSize is 2.5MB
const ChunkSize = 1024 * 1024 * 2.5

type ReadingService struct {
	logger logger.Logger
}

func NewReadingService(logger logger.Logger) *ReadingService {
	return &ReadingService{logger: logger}
}

func (r *ReadingService) Read(resource entity.Resource) chan *entity.Chunk {
	r.logger.Info("[reader]: reading started")

	chunksCh := make(chan *entity.Chunk, 1)
	go r.handleRead(resource, chunksCh)

	return chunksCh
}

func (r *ReadingService) handleRead(resource entity.Resource, chunksCh chan *entity.Chunk) {
	defer r.logger.Info("[reader]: reading stopped")
	defer close(chunksCh)

	file, err := os.Open(resource.GetPath())
	if err != nil {
		r.logger.Error(err)
		return
	}
	defer func() {
		if err = file.Close(); err != nil {
			r.logger.Error(err)
			return
		}
	}()

	for {
		chunk := entity.NewChunk(ChunkSize)

		chunk.Len, err = file.Read(chunk.Data)
		if err != nil {
			if err == io.EOF {
				r.logger.Info("[reader]: file was successfully read")
				break
			}
			r.logger.Error(err)
			return
		}

		r.sendChunk(chunk, chunksCh)
	}
}

func (r *ReadingService) sendChunk(chunk *entity.Chunk, chunksCh chan *entity.Chunk) {
	if chunk.Len == 0 {
		return
	}

	if chunk.Len < ChunkSize {
		lastChunk := make([]byte, chunk.Len)
		lastChunk = chunk.Data[:chunk.Len]
		chunk.Data = lastChunk
	}

	if chunk.Len > 0 {
		r.logger.Info(fmt.Sprintf("[reader]: read %d bytes", chunk.Len))
		chunksCh <- chunk
	}
}
