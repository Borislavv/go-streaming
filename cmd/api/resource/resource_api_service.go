package resource

import (
	"context"
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

	server := http.NewHttpServer(errCh)
	server.Listen(ctx, wg)
}
