package read

import (
	"github.com/Borislavv/video-streaming/internal/domain/model"
	"github.com/Borislavv/video-streaming/internal/domain/model/stream"
	"io"
	"log"
	"os"
)

// ChunkSize is 2.5MB
const ChunkSize = 1024 * 1024 * 2.5

type ReadingService struct {
}

func NewReadingService() *ReadingService {
	return &ReadingService{}
}

func (r *ReadingService) Read(resource model.Resource) chan *stream.Chunk {
	log.Println("[reader]: reading started")

	chunksCh := make(chan *stream.Chunk, 10)
	errCh := make(chan error)

	go r.handleReadErrs(errCh)
	go r.handleRead(resource, chunksCh, errCh)

	return chunksCh
}

func (r *ReadingService) handleRead(resource model.Resource, chunksCh chan *stream.Chunk, errCh chan error) {
	defer log.Println("[reader]: reading stopped")
	defer close(errCh)
	defer close(chunksCh)

	file, err := os.Open(resource.GetPath())
	if err != nil {
		errCh <- err
		return
	}
	defer func() {
		if err = file.Close(); err != nil {
			errCh <- err
			return
		}
	}()

	for {
		chunk := stream.NewChunk(ChunkSize)

		chunk.Len, err = file.Read(chunk.Data)
		if err != nil {
			if err == io.EOF {
				log.Println("[reader]: file was successfully read")
				break
			}
			errCh <- err
			return
		}

		if chunk.Len > 0 {
			log.Printf("[reader]: read %d bytes", chunk.Len)
			chunksCh <- chunk
		}
	}
}

func (r *ReadingService) handleReadErrs(errCh chan error) {
	for err := range errCh {
		log.Println(err)
	}
}
