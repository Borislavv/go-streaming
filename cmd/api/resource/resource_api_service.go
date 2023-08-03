package resource

import (
	"context"
	"github.com/Borislavv/video-streaming/internal/infrastructure/api/v1/controller"
	"github.com/Borislavv/video-streaming/internal/infrastructure/api/v1/controller/audio"
	"github.com/Borislavv/video-streaming/internal/infrastructure/api/v1/controller/video"
	"github.com/Borislavv/video-streaming/internal/infrastructure/server/http"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

type ResourcesApiService struct {
}

func NewApiService() *ResourcesApiService {
	return &ResourcesApiService{}
}

// Run is method which running the REST API part of app
func (r *ResourcesApiService) Run() {
	wg := &sync.WaitGroup{}
	ctx, cancel := context.WithCancel(context.Background())

	errCh := make(chan error)
	go r.handleErrors(errCh)
	defer close(errCh)

	wg.Add(1)
	go http.NewHttpServer(
		[]controller.Controller{
			video.NewCreateController(),
			video.NewDeleteVideoController(),
			video.NewGetVideoController(),
			video.NewListVideoController(),
			video.NewUpdateVideoController(),
			audio.NewCreateController(),
			audio.NewDeleteVideoController(),
			audio.NewGetVideoController(),
			audio.NewListVideoController(),
			audio.NewUpdateVideoController(),
		},
		errCh,
	).Listen(ctx, wg)
	defer func() {
		cancel()
		wg.Wait()
	}()

	stopCh := make(chan os.Signal, 1)
	signal.Notify(stopCh, os.Interrupt, syscall.SIGTERM)
	<-stopCh

	cancel()

}

// handleErrors is method which logging occurred errors
func (r *ResourcesApiService) handleErrors(errCh chan error) {
	for err := range errCh {
		log.Println(err)
	}
}
