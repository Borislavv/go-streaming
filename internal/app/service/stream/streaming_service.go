package stream

import (
	"github.com/Borislavv/video-streaming/internal/app/service"
	"github.com/Borislavv/video-streaming/internal/domain/model/video"
	"github.com/gorilla/websocket"
	"log"
	"sync"
)

const VideoPath = "/home/jared/jaredsplace/projects/go/streaming/internal/infrastructure/resource/tmp/video/example_video_new.mp4"

type StreamingService struct {
	reader service.Reader
	errCh  chan error
}

func NewStreamingService(
	reader service.Reader,
	errCh chan error,
) *StreamingService {
	return &StreamingService{
		reader: reader,
		errCh:  errCh,
	}
}

func (s *StreamingService) Stream(conn *websocket.Conn) {
	log.Println("[streamer]: start streaming")

	wg := &sync.WaitGroup{}

	wg.Add(1)
	go s.handleStream(wg, conn)

	wg.Wait()

	log.Println("[streamer]: streaming is stopped")
}

func (s *StreamingService) handleStream(wg *sync.WaitGroup, conn *websocket.Conn) {
	defer wg.Done()

	for chunk := range s.reader.Read(video.New(VideoPath)) {
		if chunk.Err != nil {
			s.errCh <- chunk.Err
			continue
		}

		if err := conn.WriteMessage(websocket.BinaryMessage, chunk.Data); err != nil {
			s.errCh <- err
			continue
		}

		log.Printf("[streamer]: wrote %d bytes to websocket\n", chunk.Len)
	}
}
