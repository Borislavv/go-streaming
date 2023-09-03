package service

import (
	"fmt"
	"github.com/Borislavv/video-streaming/internal/domain/dto"
	"github.com/Borislavv/video-streaming/internal/domain/entity"
	"io"
	"os"
)

// ChunkSize is 2.5MB
const ChunkSize = 1024 * 1024 * 2.5

type ReaderService struct {
	logger Logger
}

func NewReaderService(logger Logger) *ReaderService {
	return &ReaderService{logger: logger}
}

func (r *ReaderService) Read(resource entity.Resource) chan *dto.ChunkDto {
	r.logger.Info("[reader]: reading started")

	chunksCh := make(chan *dto.ChunkDto, 1)
	go r.handleRead(resource, chunksCh)

	return chunksCh
}

func (r *ReaderService) handleRead(resource entity.Resource, chunksCh chan *dto.ChunkDto) {
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

func (r *ReaderService) sendChunk(chunk *dto.ChunkDto, chunksCh chan *dto.ChunkDto) {
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
