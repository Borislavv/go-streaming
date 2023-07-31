package stream

import (
	"context"
	"github.com/Borislavv/video-streaming/internal/app/service/stream"
	"github.com/Borislavv/video-streaming/internal/app/service/stream/read"
	"github.com/Borislavv/video-streaming/internal/infrastructure/server/socket"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

type StreamingApiService struct {
}

func NewApiService() *StreamingApiService {
	return &StreamingApiService{}
}

// Run is method which running the streaming part of app
func (s *StreamingApiService) Run() {
	wg := &sync.WaitGroup{}
	ctx, cancel := context.WithCancel(context.Background())

	errCh := make(chan error)
	go s.handleErrors(errCh)
	defer close(errCh)

	reader := read.NewReadingService(errCh)
	streamer := stream.NewStreamingService(reader, errCh)
	server := socket.NewSocketServer(streamer, errCh)

	wg.Add(1)
	go server.Listen(ctx, wg)
	defer func() {
		cancel()
		wg.Wait()
	}()

	stopCh := make(chan os.Signal, 1)
	signal.Notify(stopCh, os.Interrupt, syscall.SIGTERM)
	<-stopCh
}

// handleErrors is method which logging occurred errors
func (s *StreamingApiService) handleErrors(errCh chan error) {
	for err := range errCh {
		log.Println(err)
	}
}
