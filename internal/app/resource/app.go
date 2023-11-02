package resource

import (
	"context"
	"github.com/Borislavv/video-streaming/internal/domain/builder"
	loggerservice "github.com/Borislavv/video-streaming/internal/domain/logger"
	"github.com/Borislavv/video-streaming/internal/domain/repository"
	authservice "github.com/Borislavv/video-streaming/internal/domain/service/authenticator"
	"github.com/Borislavv/video-streaming/internal/domain/service/extractor"
	resourceservice "github.com/Borislavv/video-streaming/internal/domain/service/resource"
	uploaderservice "github.com/Borislavv/video-streaming/internal/domain/service/uploader"
	userservice "github.com/Borislavv/video-streaming/internal/domain/service/user"
	videoservice "github.com/Borislavv/video-streaming/internal/domain/service/video"
	"github.com/Borislavv/video-streaming/internal/domain/validator"
	"github.com/Borislavv/video-streaming/internal/infrastructure/api/v1/controller"
	"github.com/Borislavv/video-streaming/internal/infrastructure/api/v1/controller/render"
	"github.com/Borislavv/video-streaming/internal/infrastructure/api/v1/controller/rest/audio"
	"github.com/Borislavv/video-streaming/internal/infrastructure/api/v1/controller/rest/auth"
	"github.com/Borislavv/video-streaming/internal/infrastructure/api/v1/controller/rest/resource"
	"github.com/Borislavv/video-streaming/internal/infrastructure/api/v1/controller/rest/user"
	"github.com/Borislavv/video-streaming/internal/infrastructure/api/v1/controller/rest/video"
	"github.com/Borislavv/video-streaming/internal/infrastructure/api/v1/controller/static"
	"github.com/Borislavv/video-streaming/internal/infrastructure/api/v1/request"
	"github.com/Borislavv/video-streaming/internal/infrastructure/api/v1/response"
	"github.com/Borislavv/video-streaming/internal/infrastructure/repository/mongodb"
	"github.com/Borislavv/video-streaming/internal/infrastructure/server/http"
	"github.com/Borislavv/video-streaming/internal/infrastructure/service/cacher"
	"github.com/Borislavv/video-streaming/internal/infrastructure/service/logger"
	"github.com/Borislavv/video-streaming/internal/infrastructure/service/storager"
	"github.com/Borislavv/video-streaming/internal/infrastructure/service/tokenizer"
	"github.com/Borislavv/video-streaming/internal/infrastructure/service/uploader"
	"github.com/Borislavv/video-streaming/internal/infrastructure/service/uploader/file"
	"github.com/caarlos0/env/v9"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
	"time"
)

const DefaultDatabaseTimeout = time.Second * 10

type ResourcesApp struct {
	cfg config
}

func NewResourcesApp() *ResourcesApp {
	return &ResourcesApp{cfg: config{}}
}

// Run is method which running the REST API part of app
func (app *ResourcesApp) Run(mWg *sync.WaitGroup) {
	defer mWg.Done()
	wg := &sync.WaitGroup{}
	ctx, cancel := context.WithCancel(context.Background())

	// init. logger and close func.
	loggerService, cls := logger.NewStdOut(ctx, app.cfg.LoggerErrorsBufferCap, app.cfg.LoggerRequestsBufferCap)
	defer func() {
		cancel()
		wg.Wait()
		cls()
	}()

	// parse env. config
	if err := env.Parse(&app.cfg); err != nil {
		loggerService.Critical(err)
		return
	}

	// connect to mongo database
	db, err := app.InitMongoDatabase(ctx, loggerService)
	if err != nil {
		loggerService.Critical(err)
		return
	}

	// request param. resolver
	reqParamsExtractor := request.NewParametersExtractor()

	// response service
	responseService := response.NewResponseService(loggerService)

	// filesystem storage
	filesystemStorage := storager.NewFilesystemStorage(ctx, loggerService)

	// filename computer
	filenameComputerService := file.NewNameService()

	// user repository
	userRepository := mongodb.NewUserRepository(db, loggerService, DefaultDatabaseTimeout)

	// user builder
	userBuilder := builder.NewUserBuilder(ctx, loggerService, reqParamsExtractor, userRepository)

	// user validator
	userValidator := validator.NewUserValidator(ctx, loggerService, userRepository, app.cfg.AdminContactEmail)

	var uploaderStrategy uploaderservice.Uploader
	if app.cfg.UploadingStrategy == uploader.MultipartFormUploadingType {
		// used parsing of full form into RAM
		uploaderStrategy =
			uploader.NewNativeUploader(
				loggerService,
				filesystemStorage,
				filenameComputerService,
				app.cfg.ResourceFormFilename,
				app.cfg.InMemoryFileSizeThreshold,
			)
	} else if app.cfg.UploadingStrategy == uploader.MultipartPartUploadingType {
		// used partial reading from multipart.Part
		uploaderStrategy =
			uploader.NewPartsUploader(
				loggerService,
				filesystemStorage,
				filenameComputerService,
			)
	}

	// auth. validator
	authValidator := validator.NewAuthValidator(loggerService, app.cfg.AdminContactEmail)

	// auth. builder
	authBuilder := builder.NewAuthBuilder(loggerService)

	// blocked token repository
	blockedTokenRepository := mongodb.NewBlockedTokenRepository(db, loggerService, DefaultDatabaseTimeout)

	/**
	 * CRUD services.
	 */

	videoBuilder, _, videoService, _ := app.InitVideoServices(
		ctx, loggerService, db, resourceValidator,
		resourceRepository, resourceService, reqParamsExtractor,
	)

	userService := userservice.NewCRUDService(
		ctx, loggerService, userBuilder, userValidator, userRepository, videoService,
	)
	tokenService := tokenizer.NewJwtService(
		ctx, loggerService, blockedTokenRepository, strings.Split(app.cfg.JwtTokenAcceptedIssuers, ","),
		app.cfg.JwtSecretSalt, app.cfg.JwtTokenIssuer, app.cfg.JwtTokenEncryptAlgo, app.cfg.JwtTokenExpiresAfter,
	)
	authService := authservice.NewAuthService(
		loggerService, userService, authValidator, tokenService,
	)

	cacheService := app.InitCacheService(ctx)

	wg.Add(1)
	go http.NewHttpServer(
		ctx,
		app.cfg.Host,
		app.cfg.Port,
		app.cfg.Transport,
		app.cfg.ApiVersionPrefix,
		app.cfg.RenderVersionPrefix,
		app.cfg.StaticVersionPrefix,
		app.InitRestApiControllers(
			cacheService,
			loggerService,
			responseService,
			resourceBuilder,
			resourceService,
			videoBuilder,
			videoService,
			userBuilder,
			userService,
			authBuilder,
			authService,
		),
		app.InitNativeRenderingControllers(
			loggerService,
			responseService,
		),
		app.InitStaticServingControllers(
			loggerService,
			responseService,
		),
		loggerService,
		reqParamsExtractor,
	).Listen(ctx, wg)

	<-app.shutdown()
}

func (app *ResourcesApp) shutdown() chan os.Signal {
	stopCh := make(chan os.Signal, 1)
	signal.Notify(stopCh, os.Interrupt, syscall.SIGTERM)
	return stopCh
}

func (app *ResourcesApp) InitCacheService(ctx context.Context) cacher.Cacher {
	s := cacher.NewMapCacheStorage(ctx)
	d := cacher.NewCacheDisplacer(ctx, time.Second*1)
	c := cacher.NewCache(s, d)
	return c
}

func (app *ResourcesApp) InitMongoDatabase(ctx context.Context, logger loggerservice.Logger) (*mongo.Database, error) {
	c, err := mongo.Connect(ctx, options.Client().ApplyURI(app.cfg.MongoUri))
	if err != nil {
		return nil, logger.CriticalPropagate(err)
	}
	defer func() { _ = c.Disconnect(ctx) }()

	if err = c.Ping(ctx, readpref.Primary()); err != nil {
		return nil, logger.CriticalPropagate(err)
	}

	return c.Database(app.cfg.MongoDb), nil
}

func (app *ResourcesApp) InitVideoServices(
	ctx context.Context,
	logger loggerservice.Logger,
	database *mongo.Database,
	resourceValidator validator.Resource,
	resourceRepository repository.Resource,
	resourceService resourceservice.CRUD,
	reqParamsExtractor extractor.RequestParams,
) (
	builder.Video,
	validator.Video,
	videoservice.CRUD,
	repository.Video,
) {
	r := mongodb.NewVideoRepository(database, logger, time.Minute)
	v := validator.NewVideoValidator(ctx, logger, resourceValidator, r, resourceRepository)
	b := builder.NewVideoBuilder(ctx, logger, reqParamsExtractor, r, resourceRepository)
	s := videoservice.NewCRUDService(ctx, logger, b, v, r, resourceService)
	return b, v, s, r
}

func (app *ResourcesApp) InitResourceServices(
	ctx context.Context,
	logger loggerservice.Logger,
	database *mongo.Database,
	uploaderStrategy uploaderservice.Uploader,
	reqParamsExtractor extractor.RequestParams,
) (
	builder.Resource,
	validator.Resource,
	resourceservice.CRUD,
	repository.Resource,
) {
	r := mongodb.NewResourceRepository(database, logger, time.Minute)
	v := validator.NewResourceValidator(ctx, r, app.cfg.MaxFilesizeThreshold)
	b := builder.NewResourceBuilder(logger, app.cfg.ResourceFormFilename, app.cfg.InMemoryFileSizeThreshold)
	resourceService := resourceservice.NewResourceService(
		ctx, logger, uploaderStrategy, v, b, r, filesystemStorage,
	)
}

func (app *ResourcesApp) InitRestApiControllers(
	cacheService cacher.Cacher,
	loggerService loggerservice.Logger,
	responseService response.Responder,
	// resource deps.
	resourceBuilder builder.Resource,
	resourceService resourceservice.CRUD,
	// video deps.
	videoBuilder builder.Video,
	videoService videoservice.CRUD,
	// user. deps.
	userBuilder builder.User,
	userService userservice.CRUD,
	// auth. deps.
	authBuilder builder.Auth,
	authService authservice.Authenticator,
) []controller.Controller {
	return []controller.Controller{
		// resource
		resource.NewUploadController(
			loggerService,
			resourceBuilder,
			resourceService,
			responseService,
		),
		// video
		video.NewCreateController(
			loggerService,
			videoBuilder,
			videoService,
			responseService,
		),
		video.NewDeleteController(
			loggerService,
			videoBuilder,
			videoService,
			responseService,
		),
		video.NewGetController(
			loggerService,
			videoBuilder,
			videoService,
			responseService,
		),
		video.NewListController(
			loggerService,
			videoBuilder,
			videoService,
			responseService,
		),
		video.NewUpdateController(
			loggerService,
			videoBuilder,
			videoService,
			responseService,
		),
		// audio
		audio.NewCreateController(),
		audio.NewDeleteController(),
		audio.NewGetController(),
		audio.NewListController(),
		audio.NewUpdateController(),
		// user
		user.NewCreateController(
			loggerService,
			userBuilder,
			userService,
			responseService,
		),
		user.NewUpdateUserController(
			loggerService,
			userBuilder,
			userService,
			responseService,
		),
		user.NewGetController(
			loggerService,
			userBuilder,
			userService,
			cacheService,
			responseService,
		),
		user.NewDeleteController(
			loggerService,
			userBuilder,
			userService,
			responseService,
		),
		// auth
		auth.NewAuthorizationController(
			loggerService,
			authBuilder,
			authService,
			responseService,
		),
		auth.NewRegistrationController(
			loggerService,
			userBuilder,
			userService,
			responseService,
		),
	}
}

func (app *ResourcesApp) InitNativeRenderingControllers(
	loggerService loggerservice.Logger,
	responseService response.Responder,
) []controller.Controller {
	return []controller.Controller{
		render.NewIndexController(loggerService, responseService),
	}
}

func (app *ResourcesApp) InitStaticServingControllers(
	loggerService loggerservice.Logger,
	responseService response.Responder,
) []controller.Controller {
	return []controller.Controller{
		static.NewFilesController(loggerService, responseService),
	}
}
