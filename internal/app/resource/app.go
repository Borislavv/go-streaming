package resource

import (
	"context"
	"github.com/Borislavv/video-streaming/internal/domain/builder"
	"github.com/Borislavv/video-streaming/internal/domain/service"
	"github.com/Borislavv/video-streaming/internal/domain/validator"
	"github.com/Borislavv/video-streaming/internal/infrastructure/api/v1/controller"
	"github.com/Borislavv/video-streaming/internal/infrastructure/api/v1/controller/render"
	"github.com/Borislavv/video-streaming/internal/infrastructure/api/v1/controller/rest/audio"
	"github.com/Borislavv/video-streaming/internal/infrastructure/api/v1/controller/rest/video"
	"github.com/Borislavv/video-streaming/internal/infrastructure/api/v1/controller/static"
	"github.com/Borislavv/video-streaming/internal/infrastructure/api/v1/request"
	"github.com/Borislavv/video-streaming/internal/infrastructure/api/v1/response"
	"github.com/Borislavv/video-streaming/internal/infrastructure/logger"
	"github.com/Borislavv/video-streaming/internal/infrastructure/repository/mongodb"
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

type ResourcesApp struct {
	cfg config
}

func NewResourcesApp() *ResourcesApp {
	return &ResourcesApp{cfg: config{}}
}

// Run is method which running the REST API part of app
func (r *ResourcesApp) Run(mWg *sync.WaitGroup) {
	defer mWg.Done()
	wg := &sync.WaitGroup{}
	ctx, cancel := context.WithCancel(context.Background())

	// init. loggerService and close func.
	loggerService, cls := logger.NewCliLogger(1)
	defer func() {
		cancel()
		wg.Wait()
		cls()
	}()

	// parse env. config
	if err := env.Parse(&r.cfg); err != nil {
		loggerService.Critical(err)
		return
	}

	// init. mongodb client
	mongoClient, err := mongo.Connect(ctx, options.Client().ApplyURI(r.cfg.MongoUri))
	if err != nil {
		loggerService.Critical(err)
		return
	}
	defer mongoClient.Disconnect(ctx)

	// ping mongodb
	if err = mongoClient.Ping(ctx, readpref.Primary()); err != nil {
		loggerService.Critical(err)
		return
	}

	// connect to target mongodb database
	db := mongoClient.Database(r.cfg.MongoDb)

	reqParamsExtractor := request.NewParametersExtractor()

	// init. video repository
	videoRepository := mongodb.NewVideoRepository(db, loggerService, time.Minute)

	// init. video validator
	videoValidator := validator.NewVideoValidator(ctx, videoRepository)

	// init. video builder
	videoBuilder := builder.NewVideoBuilder(ctx, reqParamsExtractor, videoRepository)

	// init. video service
	videoService := service.NewVideoService(ctx, loggerService, videoBuilder, videoValidator, videoRepository)

	// init. response service
	responseService := response.NewResponseService(loggerService)

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
				videoBuilder,
				videoService,
				responseService,
			),
			video.NewDeleteVideoController(
				videoBuilder,
				videoService,
				responseService,
			),
			video.NewGetVideoController(
				videoBuilder,
				videoService,
				responseService,
			),
			video.NewListVideoController(
				videoBuilder,
				videoService,
				responseService,
			),
			video.NewUpdateVideoController(
				videoBuilder,
				videoService,
				responseService,
			),
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
		loggerService,
	).Listen(ctx, wg)

	stopCh := make(chan os.Signal, 1)
	signal.Notify(stopCh, os.Interrupt, syscall.SIGTERM)
	<-stopCh
}
