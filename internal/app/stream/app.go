package stream

import (
	"context"
	"github.com/Borislavv/video-streaming/internal/domain/service"
	"github.com/Borislavv/video-streaming/internal/infrastructure/logger"
	"github.com/Borislavv/video-streaming/internal/infrastructure/server/socket"
	"github.com/caarlos0/env/v9"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

type StreamingApp struct {
	cfg config
}

func NewStreamingApp() *StreamingApp {
	return &StreamingApp{cfg: config{}}
}

// Run is method which running the streaming part of app
func (s *StreamingApp) Run(mWg *sync.WaitGroup) {
	defer mWg.Done()
	wg := &sync.WaitGroup{}
	ctx, cancel := context.WithCancel(context.Background())

	// init. logger and close func.
	loggerService, cls := logger.NewCliLogger(1)
	defer func() {
		cancel()
		wg.Wait()
		cls()
	}()

	if err := env.Parse(&s.cfg); err != nil {
		loggerService.Critical(err)
		return
	}

	readerService := service.NewReaderService(loggerService)
	streamingService := service.NewStreamingService(readerService, loggerService)
	server := socket.NewSocketServer(s.cfg.Host, s.cfg.Port, s.cfg.Transport, streamingService, loggerService)

	wg.Add(1)
	go server.Listen(ctx, wg)

	stopCh := make(chan os.Signal, 1)
	signal.Notify(stopCh, os.Interrupt, syscall.SIGTERM)
	<-stopCh
}
