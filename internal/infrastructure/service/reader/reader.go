package reader

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"github.com/Borislavv/video-streaming/internal/domain/dto"
	"github.com/Borislavv/video-streaming/internal/domain/logger"
	"io"
	"os"
)

const (
	ChunkSize    = 1024 * 1024 * 2.5 // 2.5MB
	ChunksBuffer = 4
)

type ResourceReader struct {
	logger logger.Logger
	cache  map[string][]dto.Chunk
}

func NewReaderService(logger logger.Logger) *ResourceReader {
	return &ResourceReader{
		logger: logger,
		cache:  map[string][]dto.Chunk{},
	}
}

// Read will read a resource and send file as butches of bytes
func (r *ResourceReader) Read(resource dto.Resource) chan dto.Chunk {
	r.logger.Info(fmt.Sprintf("recourse '%v' reading started", resource.GetFilepath()))

	chunksCh := make(chan dto.Chunk, ChunksBuffer)
	go r.handleRead(resource, chunksCh)

	return chunksCh
}

func (r *ResourceReader) handleRead(resource dto.Resource, chunksCh chan dto.Chunk) {
	defer func() {
		close(chunksCh)
		r.logger.Info(fmt.Sprintf("recourse '%v' reading finished", resource.GetFilepath()))
	}()

	file, err := os.Open(resource.GetFilepath())
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
				break
			}
			r.logger.Error(err)
			return
		}

		r.sendChunk(chunk, chunksCh)
	}
}

func (r *ResourceReader) sendChunk(chunk *dto.ChunkDTO, chunksCh chan dto.Chunk) {
	if chunk.Len == 0 {
		return
	}

	if chunk.Len < ChunkSize {
		lastChunk := make([]byte, chunk.Len)
		lastChunk = chunk.Data[:chunk.Len]
		chunk.Data = lastChunk
	}

	if chunk.Len > 0 {
		r.logger.Info(fmt.Sprintf("%d bytes read and sent", chunk.Len))
		chunksCh <- chunk
	}
}

func (r *ResourceReader) cached(resource dto.Resource) ([]dto.Chunk, error) {
	hash := md5.New()
	if _, err := hash.Write([]byte(resource.GetFilepath())); err != nil {
		r.logger.Emergency(err)
		return nil, err
	}
	key := hex.EncodeToString(hash.Sum(nil))

	if data, found := r.cache[key]; found {
		return data, nil
	}

	r.cache[key] =

	return nil, nil
}
