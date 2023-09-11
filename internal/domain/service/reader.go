package service

import (
	"fmt"
	"github.com/Borislavv/video-streaming/internal/domain/dto"
	"github.com/Borislavv/video-streaming/internal/domain/logger"
	"io"
	"os"
)

// ChunkSize is 2.5MB
const ChunkSize = 1024 * 1024 * 2.5

type ResourceReader struct {
	logger logger.Logger
}

func NewReaderService(logger logger.Logger) *ResourceReader {
	return &ResourceReader{logger: logger}
}

func (r *ResourceReader) Read(resource dto.Resource) chan *dto.Chunk {
	r.logger.Info("[reader]: reading started")

	chunksCh := make(chan *dto.Chunk, 1)
	go r.handleRead(resource, chunksCh)

	return chunksCh
}

func (r *ResourceReader) handleRead(resource dto.Resource, chunksCh chan *dto.Chunk) {
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
		chunk := dto.NewChunk(ChunkSize)

		chunk.Len, err = file.Read(chunk.Data)
		if err != nil {
			if err == io.EOF {
				r.logger.Info("[reader]: file was successfully reader")
				break
			}
			r.logger.Error(err)
			return
		}

		r.sendChunk(chunk, chunksCh)
	}
}

func (r *ResourceReader) sendChunk(chunk *dto.Chunk, chunksCh chan *dto.Chunk) {
	if chunk.Len == 0 {
		return
	}

	if chunk.Len < ChunkSize {
		lastChunk := make([]byte, chunk.Len)
		lastChunk = chunk.Data[:chunk.Len]
		chunk.Data = lastChunk
	}

	if chunk.Len > 0 {
		r.logger.Info(fmt.Sprintf("[reader]: reader %d bytes", chunk.Len))
		chunksCh <- chunk
	}
}
