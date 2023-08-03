package resource

import (
	"context"
	"github.com/Borislavv/video-streaming/internal/infrastructure/api/v1/controller"
	"github.com/Borislavv/video-streaming/internal/infrastructure/api/v1/controller/audio"
	"github.com/Borislavv/video-streaming/internal/infrastructure/api/v1/controller/video"
	"github.com/Borislavv/video-streaming/internal/infrastructure/server/http"
	"sync"
)

type ResourcesApiService struct {
}

func NewApiService() *ResourcesApiService {
	return &ResourcesApiService{}
}

func (r *ResourcesApiService) Run() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	wg := &sync.WaitGroup{}
	errCh := make(chan error)

	http.NewHttpServer(
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
}
