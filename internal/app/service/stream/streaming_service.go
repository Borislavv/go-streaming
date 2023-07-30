package stream

import (
	"github.com/Borislavv/video-streaming/internal/app/service"
	"github.com/Borislavv/video-streaming/internal/domain/model/video"
	"github.com/gorilla/websocket"
	"log"
	"sync"
)

const VideoPath = "/home/jared/jaredsplace/projects/go/streaming/internal/infrastructure/resource/tmp/video/example_video.mp4"

type StreamingService struct {
	reader service.Reader
}

func NewStreamingService(
	reader service.Reader,
) *StreamingService {
	return &StreamingService{
		reader: reader,
	}
}

func (s *StreamingService) Stream(conn *websocket.Conn) {
	log.Println("streaming service: started")

	wg := &sync.WaitGroup{}
	errCh := make(chan error)

	wg.Add(1)
	go s.handleStreamErrs(wg, errCh)

	wg.Add(1)
	go s.handleStream(wg, conn, errCh)

	wg.Wait()

	log.Println("streaming service: stopped")
}

func (s *StreamingService) handleStream(wg *sync.WaitGroup, conn *websocket.Conn, errCh chan error) {
	defer wg.Done()
	defer close(errCh)

	for chunk := range s.reader.Read(video.New(VideoPath)) {
		if chunk.Err != nil {
			errCh <- chunk.Err
			continue
		}

		if err := conn.WriteMessage(websocket.BinaryMessage, chunk.Data); err != nil {
			errCh <- err
			continue
		}

		log.Printf("streaming service: wrote %d bytes to websocket\n", chunk.Len)
	}
}

func (s *StreamingService) handleStreamErrs(wg *sync.WaitGroup, errCh chan error) {
	defer wg.Done()
	for err := range errCh {
		log.Println(err)
	}
}
