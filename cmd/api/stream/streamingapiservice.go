package stream

import (
	"context"
	"github.com/Borislavv/video-streaming/internal/app/service/stream"
	"github.com/Borislavv/video-streaming/internal/app/service/stream/read"
	"github.com/Borislavv/video-streaming/internal/infrastructure/logger/cli"
	"github.com/Borislavv/video-streaming/internal/infrastructure/server/socket"
	"github.com/caarlos0/env/v9"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

type config struct {
	// server
	Host      string `env:"RESOURCES_SERVER_HOST" envDefault:"0.0.0.0"`
	Port      string `env:"RESOURCES_SERVER_PORT" envDefault:"9988"`
	Transport string `env:"RESOURCES_SERVER_TRANSPORT_PROTOCOL" envDefault:"tcp"`
	// database
	MongoUri string `env:"MONGO_URI" envDefault:"mongodb://database:27017/streaming"`
}

type StreamingApiService struct {
	cfg config
}

func NewApiService() *StreamingApiService {
	return &StreamingApiService{cfg: config{}}
}

// Run is method which running the streaming part of app
func (s *StreamingApiService) Run(mWg *sync.WaitGroup) {
	defer mWg.Done()
	wg := &sync.WaitGroup{}
	ctx, cancel := context.WithCancel(context.Background())

	errCh := make(chan error, 1)
	logger := cli.NewLogger(errCh)
	defer func() {
		cancel()
		wg.Wait()
		close(errCh)
	}()

	if err := env.Parse(&s.cfg); err != nil {
		logger.Critical(err)

	}

	reader := read.NewReadingService(logger)
	streamer := stream.NewStreamingService(reader, logger)
	server := socket.NewSocketServer(s.cfg.Host, s.cfg.Port, s.cfg.Transport, streamer, logger)

	wg.Add(1)
	go server.Listen(ctx, wg)

	stopCh := make(chan os.Signal, 1)
	signal.Notify(stopCh, os.Interrupt, syscall.SIGTERM)
	<-stopCh
}
