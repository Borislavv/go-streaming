package video

import (
	"github.com/Borislavv/video-streaming/internal/domain/model/stream"
	"io"
	"log"
	"os"
)

// ChunkSize is 2.5MB
const ChunkSize = 1024 * 1024 * 2.5

type ManagerService struct {
}

func NewVideoManagerService() *ManagerService {
	return &ManagerService{}
}

func (m *ManagerService) Read() chan *stream.Chunk {
	log.Println("video manager: reading started")

	chunksCh := make(chan *stream.Chunk, 100) // 250mb chunks buffer, each chunk by 2.5mb

	go func() {
		defer close(chunksCh)
		defer log.Println("video manager: reading stopped")

		file, err := os.Open("/home/jared/jaredsplace/projects/go/streaming/internal/infrastructure/resource/tmp/video/example_video.mp4")
		if err != nil {
			log.Fatalln(err)
		}
		defer file.Close()

		for {
			chunk := &stream.Chunk{Data: make([]byte, ChunkSize)}

			chunk.Len, err = file.Read(chunk.Data)
			if err != nil {
				if err != io.EOF {
					chunk.Err = err
					chunksCh <- chunk
					return
				}
				log.Println("video manager: file was fully read")
				break
			}

			log.Printf("video manager: readed %d bytes\n", chunk.Len)

			chunksCh <- chunk
		}
	}()

	return chunksCh
}
