package video

import (
	"github.com/Borislavv/video-streaming/internal/domain/model"
	"io"
	"log"
	"os"
)

const ChunkSize = 1024 * 1024

type ManagerService struct {
}

func NewVideoManagerService() *ManagerService {
	return &ManagerService{}
}

func (m *ManagerService) Read() chan *model.Chunk {
	log.Println("video manager: reading started")

	chunksCh := make(chan *model.Chunk) // 10mb chunks buffer, each chunk by 1mb

	go func() {
		defer close(chunksCh)
		defer log.Println("video manager: reading stopped")

		file, err := os.Open("/home/jared/jaredsplace/projects/go/streaming/internal/infrastructure/resource/tmp/video/example_video.mp4")
		if err != nil {
			log.Fatalln(err)
		}
		defer file.Close()

		for {
			chunk := &model.Chunk{Data: make([]byte, ChunkSize)}

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
