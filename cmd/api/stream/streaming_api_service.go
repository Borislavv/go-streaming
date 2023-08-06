package stream

import (
	"context"
	"github.com/Borislavv/video-streaming/internal/app/service/stream"
	"github.com/Borislavv/video-streaming/internal/app/service/stream/read"
	"github.com/Borislavv/video-streaming/internal/infrastructure/logger/cli"
	"github.com/Borislavv/video-streaming/internal/infrastructure/server/socket"
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
func (s *StreamingApiService) Run(mWg *sync.WaitGroup) {
	defer mWg.Done()
	wg := &sync.WaitGroup{}
	ctx, cancel := context.WithCancel(context.Background())

	errCh := make(chan error, 1)
	logger := cli.NewLogger(errCh)
	defer close(errCh)

	reader := read.NewReadingService(logger)
	streamer := stream.NewStreamingService(reader, logger)
	server := socket.NewSocketServer(streamer, logger)

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
