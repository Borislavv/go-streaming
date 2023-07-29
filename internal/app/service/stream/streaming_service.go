package stream

import (
	"github.com/Borislavv/video-streaming/internal/app/service"
	"github.com/gorilla/websocket"
	"log"
)

type StreamingService struct {
	manager service.Manager
}

func NewStreamingService(manager service.Manager) *StreamingService {
	return &StreamingService{
		manager: manager,
	}
}

func (s *StreamingService) StartStreaming(conn *websocket.Conn) error {
	defer conn.Close()

	log.Println("streaming service: started")

	chunkCh := s.manager.Read()

	for chunk := range chunkCh {
		if chunk.Err != nil {
			return chunk.Err
		}

		err := conn.WriteMessage(websocket.BinaryMessage, chunk.Data)
		if err != nil {
			return err
		}

		log.Printf("streaming service: wrote message to websocket")
	}

	log.Println("streaming service: stopped")

	return nil
}
