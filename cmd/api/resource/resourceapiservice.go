package resource

import (
	"context"
	aggbuilder "github.com/Borislavv/video-streaming/internal/domain/builder/agg"
	dtobuilder "github.com/Borislavv/video-streaming/internal/domain/builder/dto"
	entitybuilder "github.com/Borislavv/video-streaming/internal/domain/builder/entity"
	"github.com/Borislavv/video-streaming/internal/infrastructure/api/v1/controller"
	"github.com/Borislavv/video-streaming/internal/infrastructure/api/v1/controller/render"
	"github.com/Borislavv/video-streaming/internal/infrastructure/api/v1/controller/rest/audio"
	"github.com/Borislavv/video-streaming/internal/infrastructure/api/v1/controller/rest/video"
	"github.com/Borislavv/video-streaming/internal/infrastructure/api/v1/controller/static"
	"github.com/Borislavv/video-streaming/internal/infrastructure/logger/cli"
	mongorepository "github.com/Borislavv/video-streaming/internal/infrastructure/repository/mongo"
	"github.com/Borislavv/video-streaming/internal/infrastructure/server/http"
	"github.com/caarlos0/env/v9"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
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
	MongoDb  string `env:"MONGO_DATABASE" envDefault:"streaming"`
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

	// init. logger
	errCh := make(chan error, 1)
	logger := cli.NewLogger(errCh)
	defer func() {
		cancel()
		wg.Wait()
		close(errCh)
	}()

	// parse env. config
	if err := env.Parse(&r.cfg); err != nil {
		logger.Critical(err)
		return
	}

	// init. mongodb client
	mongoClient, err := mongo.Connect(ctx, options.Client().ApplyURI(r.cfg.MongoUri))
	if err != nil {
		logger.Critical(err)
		return
	}
	defer mongoClient.Disconnect(ctx)

	// ping mongodb
	if err = mongoClient.Ping(ctx, readpref.Primary()); err != nil {
		logger.Critical(err)
		return
	}

	// connect to target mongo database
	mongodb := mongoClient.Database(r.cfg.MongoDb)

	// create video repository
	videoRepository := mongorepository.NewVideoRepository(mongodb, time.Minute)

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
			video.NewCreateController(
				logger,
				dtobuilder.NewVideoDtoBuilder(),
				entitybuilder.NewVideoEntityBuilder(),
				aggbuilder.NewVideoAggBuilder(),
				videoRepository,
			),
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
