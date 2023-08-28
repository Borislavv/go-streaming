package resource

import (
	"context"
	"github.com/Borislavv/video-streaming/internal/infrastructure/api/v1/controller"
	"github.com/Borislavv/video-streaming/internal/infrastructure/api/v1/controller/render"
	"github.com/Borislavv/video-streaming/internal/infrastructure/api/v1/controller/rest/audio"
	"github.com/Borislavv/video-streaming/internal/infrastructure/api/v1/controller/rest/video"
	"github.com/Borislavv/video-streaming/internal/infrastructure/api/v1/controller/static"
	"github.com/Borislavv/video-streaming/internal/infrastructure/logger/cli"
	"github.com/Borislavv/video-streaming/internal/infrastructure/server/http"
	"github.com/caarlos0/env/v9"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

const (
	MongoVideoCollection = "videos"
	MongoAudioCollection = "audio"
)

type config struct {
	// api
	ApiVersionPrefix    string `env:"API_VERSION_PREFIX" envDefault:"/api/v1"`
	RenderVersionPrefix string `env:"RENDER_VERSION_PREFIX" envDefault:""`
	StaticVersionPrefix string `env:"STATIC_VERSION_PREFIX" envDefault:""`
	// server
	Host      string `env:"RESOURCES_SERVER_HOST" envDefault:"0.0.0.0"`
	Port      string `env:"RESOURCES_SERVER_PORT" envDefault:"8000"`
	Transport string `env:"RESOURCES_SERVER_TRANSPORT_PROTOCOL" envDefault:"tcp"`
	// database
	MongoUri string `env:"MONGO_URI" envDefault:"mongodb://database:27017/streaming"`
}

type ResourcesApiService struct {
	cfg config
}

func NewApiService() *ResourcesApiService {
	return &ResourcesApiService{cfg: config{}}
}

// Run is method which running the REST API part of app
func (r *ResourcesApiService) Run(mWg *sync.WaitGroup) {
	defer mWg.Done()
	wg := &sync.WaitGroup{}
	ctx, cancel := context.WithCancel(context.Background())
	defer func() {
		cancel()
		wg.Wait()
	}()

	errCh := make(chan error, 1)
	logger := cli.NewLogger(errCh)
	defer close(errCh)

	if err := env.Parse(&r.cfg); err != nil {
		logger.Critical(err)
		return
	}

	wg.Add(1)
	go http.NewHttpServer(
		r.cfg.Host,
		r.cfg.Port,
		r.cfg.Transport,
		r.cfg.ApiVersionPrefix,
		r.cfg.RenderVersionPrefix,
		r.cfg.StaticVersionPrefix,
		// rest api controllers
		[]controller.Controller{
			// video
			video.NewCreateController(),
			video.NewDeleteVideoController(),
			video.NewGetVideoController(),
			video.NewListVideoController(),
			video.NewUpdateVideoController(),
			// audio
			audio.NewCreateController(),
			audio.NewDeleteVideoController(),
			audio.NewGetVideoController(),
			audio.NewListVideoController(),
			audio.NewUpdateVideoController(),
		},
		// native rendering controllers
		[]controller.Controller{
			render.NewIndexController(),
		},
		// static serving controllers
		[]controller.Controller{
			static.NewResourceController(),
		},
		logger,
	).Listen(ctx, wg)

	stopCh := make(chan os.Signal, 1)
	signal.Notify(stopCh, os.Interrupt, syscall.SIGTERM)
	<-stopCh
}
