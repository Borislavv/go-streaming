package file

import (
	"github.com/Borislavv/video-streaming/internal/domain/model/file"
	"io"
	"log"
	"os"
)

const DefaultChunkSize = uint64(1024 * 1024 * 1)

type MTReader struct {
}

func NewMTReader() *MTReader {
	return &MTReader{}
}

func (m *MTReader) Read(path string, chunkSizeCh <-chan uint64, resCh chan<- *file.Chunk, errCh chan<- error) {
	f, err := os.Open(path)
	if err != nil {
		errCh <- err
		return
	}
	defer func() {
		if err = f.Close(); err != nil {
			errCh <- err
			return
		}
	}()

	go func() {
		defer close(errCh)
		defer close(resCh)

		chunkSize := DefaultChunkSize

		for {
			select {
			case chunkSize = <-chunkSizeCh:
			default:
				m.readChunk(f, chunkSize, resCh, errCh)
			}
		}
	}()
}

func (m *MTReader) readChunk(f *os.File, chunkSize uint64, resCh chan<- *file.Chunk, errCh chan<- error) {
	chunk := file.NewChunk(chunkSize)
	n, err := f.Read(chunk.Bytes)
	if err != nil {
		if err == io.EOF {
			log.Println("[MTReader]: file was successfully read")
			return
		}
		errCh <- err
		return
	}
	chunk.Len = n
	resCh <- chunk
	log.Printf("[reader]: read %d bytes\n", chunk.Len)
}
