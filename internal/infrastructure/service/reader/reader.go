package reader

import (
	"fmt"
	"github.com/Borislavv/video-streaming/internal/domain/dto"
	"github.com/Borislavv/video-streaming/internal/domain/logger"
	"io"
	"os"
)

const (
	ChunkSize      = 1024 * 1024 * 1 // 1MB
	ChunksChBuffer = 1
)

type ResourceReader struct {
	logger logger.Logger
	cache  map[string]map[int]dto.Chunk
}

func NewReaderService(logger logger.Logger) *ResourceReader {
	return &ResourceReader{
		logger: logger,
		cache:  map[string]map[int]dto.Chunk{},
	}
}

// Read will read a resource and send file as butches of bytes
func (r *ResourceReader) Read(resource dto.Resource) chan dto.Chunk {
	r.logger.Info(fmt.Sprintf("recourse '%v' reading started", resource.GetName()))

	chunksCh := make(chan dto.Chunk, ChunksChBuffer)
	go r.handleRead(resource, chunksCh)

	return chunksCh
}

func (r *ResourceReader) handleRead(resource dto.Resource, chunksCh chan dto.Chunk) {
	defer func() {
		close(chunksCh)

	}()

	file, ferr := os.Open(resource.GetFilepath())
	if ferr != nil {
		r.logger.Error(ferr)
		return
	}
	defer func() { _ = file.Close() }()

	for {
		chunk := dto.NewChunk(ChunkSize, 0)

		length, err := file.Read(chunk.Data)
		if err != nil {
			if err == io.EOF {
				break
			}
			r.logger.Error(err)
			break
		}
		chunk.SetLen(length)

		if chunk.GetLen() < ChunkSize {
			chunk.SetData(chunk.GetData()[:chunk.GetLen()])
		}

		if chunk.GetLen() > 0 {
			chunksCh <- chunk
		}
	}

	r.logger.Info(fmt.Sprintf("recourse '%v' reading finished", resource.GetName()))
}
