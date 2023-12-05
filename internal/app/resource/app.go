package resource

import (
	"context"
	"github.com/Borislavv/video-streaming/internal/domain/builder"
	loggerservice "github.com/Borislavv/video-streaming/internal/domain/logger"
	"github.com/Borislavv/video-streaming/internal/domain/repository"
	"github.com/Borislavv/video-streaming/internal/domain/service/accessor"
	authservice "github.com/Borislavv/video-streaming/internal/domain/service/authenticator"
	cacheservice "github.com/Borislavv/video-streaming/internal/domain/service/cacher"
	"github.com/Borislavv/video-streaming/internal/domain/service/extractor"
	resourceservice "github.com/Borislavv/video-streaming/internal/domain/service/resource"
	securityservice "github.com/Borislavv/video-streaming/internal/domain/service/security"
	storagerservice "github.com/Borislavv/video-streaming/internal/domain/service/storager"
	tokenizerservice "github.com/Borislavv/video-streaming/internal/domain/service/tokenizer"
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
	"github.com/Borislavv/video-streaming/internal/infrastructure/repository/storage/cache"
	"github.com/Borislavv/video-streaming/internal/infrastructure/repository/storage/mongodb"
	"github.com/Borislavv/video-streaming/internal/infrastructure/server/http"
	"github.com/Borislavv/video-streaming/internal/infrastructure/service/cacher"
	"github.com/Borislavv/video-streaming/internal/infrastructure/service/logger"
	"github.com/Borislavv/video-streaming/internal/infrastructure/service/security"
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
	client, db, err := app.InitMongoDatabase(ctx, loggerService)
	if err != nil {
		loggerService.Critical(err)
		return
	}
	defer func() { _ = client.Disconnect(ctx) }()

	// Cache dependencies initialization
	cacheService := app.InitCacheService(ctx)

	// Request-Response dependencies initialization
	requestService, responseService := app.InitRequestResponseServices(loggerService)

	// Access service
	accessService := accessor.NewAccessService(loggerService)

	// Files uploader dependencies initialization
	uploadingStorage, _, uploadingStrategy := app.InitUploaderServices(
		ctx, loggerService,
	)

	// Resource dependencies initialization
	resourceBuilder, resourceValidator, resourceService, resourceRepository := app.InitResourceServices(
		ctx, loggerService, db, uploadingStrategy, uploadingStorage, cacheService,
	)

	// Video dependencies initialization
	videoBuilder, _, videoService, _ := app.InitVideoServices(
		ctx, loggerService, db, resourceValidator, resourceRepository,
		resourceService, requestService, accessService, cacheService,
	)

	passwordService := app.InitPasswordService(loggerService)

	// User dependencies initialization
	userBuilder, _, userService, _ := app.InitUserServices(
		ctx, loggerService, db, videoService, requestService, passwordService, cacheService,
	)

	// Token dependencies initialization
	_, tokenService := app.InitTokenServices(
		ctx, loggerService, db,
	)

	// Auth dependencies initialization
	authBuilder, _, authService := app.InitAuthServices(
		loggerService, tokenService, userService, passwordService,
	)

	wg.Add(1)
	go http.NewHttpServer(
		ctx,
		app.cfg.Host,
		app.cfg.Port,
		app.cfg.Transport,
		app.cfg.ApiVersionPrefix,
		app.cfg.RenderVersionPrefix,
		app.cfg.StaticVersionPrefix,
		app.InitAuthedRestApiControllers(
			cacheService,
			loggerService,
			responseService,
			resourceBuilder,
			resourceService,
			videoBuilder,
			videoService,
			userBuilder,
			userService,
			authService,
		),
		app.InitUnauthedRestApiControllers(
			loggerService,
			responseService,
			userBuilder,
			userService,
			authBuilder,
			authService,
		),
		app.InitAuthedNativeRenderingControllers(
			loggerService,
			responseService,
		),
		app.InitUnauthedNativeRenderingControllers(
			loggerService,
			responseService,
		),
		app.InitStaticServingControllers(
			loggerService,
			responseService,
		),
		loggerService,
		authService,
		requestService,
		responseService,
	).Listen(ctx, wg)

	<-app.shutdown()
}

func (app *ResourcesApp) shutdown() chan os.Signal {
	stopCh := make(chan os.Signal, 1)
	signal.Notify(stopCh, os.Interrupt, syscall.SIGTERM)
	return stopCh
}

func (app *ResourcesApp) InitPasswordService(logger loggerservice.Logger) securityservice.PasswordHasher {
	return security.NewPasswordHasher(logger, 14)
}

func (app *ResourcesApp) InitCacheService(ctx context.Context) cacheservice.Cacher {
	s := cacher.NewMapCacheStorage(ctx)
	d := cacher.NewCacheDisplacer(ctx, time.Second*1)
	c := cacher.NewCache(s, d)
	return c
}

func (app *ResourcesApp) InitMongoDatabase(
	ctx context.Context,
	logger loggerservice.Logger,
) (
	*mongo.Client,
	*mongo.Database,
	error,
) {
	c, err := mongo.Connect(ctx, options.Client().ApplyURI(app.cfg.MongoUri))
	if err != nil {
		return nil, nil, logger.CriticalPropagate(err)
	}

	if err = c.Ping(ctx, readpref.Primary()); err != nil {
		return nil, nil, logger.CriticalPropagate(err)
	}

	return c, c.Database(app.cfg.MongoDb), nil
}

func (app *ResourcesApp) InitVideoServices(
	ctx context.Context,
	logger loggerservice.Logger,
	database *mongo.Database,
	resourceValidator validator.Resource,
	resourceRepository repository.Resource,
	resourceService resourceservice.CRUD,
	reqParamsExtractor extractor.RequestParams,
	accessService accessor.Accessor,
	cacher cacheservice.Cacher,
) (
	builder.Video,
	validator.Video,
	videoservice.CRUD,
	repository.Video,
) {
	r := mongodb.NewVideoRepository(database, logger, time.Minute)
	c := cache.NewVideoRepository(logger, cacher, r)
	v := validator.NewVideoValidator(ctx, logger, resourceValidator, accessService, r, resourceRepository)
	b := builder.NewVideoBuilder(ctx, logger, reqParamsExtractor, r, resourceRepository)
	s := videoservice.NewCRUDService(ctx, logger, b, v, r, resourceService)
	return b, v, s, c
}

func (app *ResourcesApp) InitResourceServices(
	ctx context.Context,
	logger loggerservice.Logger,
	database *mongo.Database,
	uploader uploaderservice.Uploader,
	storage storagerservice.Storage,
	cacher cacheservice.Cacher,
) (
	builder.Resource,
	validator.Resource,
	resourceservice.CRUD,
	repository.Resource,
) {
	r := mongodb.NewResourceRepository(database, logger, time.Minute)
	c := cache.NewResourceRepository(logger, cacher, r)
	v := validator.NewResourceValidator(ctx, r, app.cfg.MaxFilesizeThreshold)
	b := builder.NewResourceBuilder(logger, app.cfg.ResourceFormFilename, app.cfg.InMemoryFileSizeThreshold)
	s := resourceservice.NewResourceService(ctx, logger, uploader, v, b, r, storage)
	return b, v, s, c
}

func (app *ResourcesApp) InitUserServices(
	ctx context.Context,
	logger loggerservice.Logger,
	database *mongo.Database,
	videoService videoservice.CRUD,
	reqParamsExtractor extractor.RequestParams,
	passwordHasher securityservice.PasswordHasher,
	cacher cacheservice.Cacher,
) (
	builder.User,
	validator.User,
	userservice.CRUD,
	repository.User,
) {
	r := mongodb.NewUserRepository(database, logger, DefaultDatabaseTimeout)
	c := cache.NewUserRepository(logger, cacher, r)
	b := builder.NewUserBuilder(ctx, logger, reqParamsExtractor, r, passwordHasher)
	v := validator.NewUserValidator(ctx, logger, r, app.cfg.AdminContactEmail)
	s := userservice.NewCRUDService(ctx, logger, b, v, r, videoService)
	return b, v, s, c
}

func (app *ResourcesApp) InitAuthServices(
	logger loggerservice.Logger,
	tokenService tokenizerservice.Tokenizer,
	userService userservice.CRUD,
	passwordHasher securityservice.PasswordHasher,
) (
	builder.Auth,
	validator.Auth,
	authservice.Authenticator,
) {
	v := validator.NewAuthValidator(logger, app.cfg.AdminContactEmail)
	b := builder.NewAuthBuilder(logger)
	s := authservice.NewAuthService(logger, userService, v, tokenService, passwordHasher)
	return b, v, s
}

func (app *ResourcesApp) InitTokenServices(
	ctx context.Context,
	logger loggerservice.Logger,
	database *mongo.Database,
) (
	repository.BlockedToken,
	tokenizerservice.Tokenizer,
) {
	r := mongodb.NewBlockedTokenRepository(database, logger, DefaultDatabaseTimeout)
	s := tokenizer.NewJwtService(
		ctx, logger, r, strings.Split(app.cfg.JwtTokenAcceptedIssuers, ","),
		app.cfg.JwtSecretSalt, app.cfg.JwtTokenIssuer, app.cfg.JwtTokenEncryptAlgo, app.cfg.JwtTokenExpiresAfter,
	)
	return r, s
}

func (app *ResourcesApp) InitUploaderServices(
	ctx context.Context,
	logger loggerservice.Logger,
) (
	storagerservice.Storage,
	file.NameComputer,
	uploaderservice.Uploader,
) {
	// filesystem storage
	filesystemStorage := storager.NewFilesystemStorage(ctx, logger)

	// filename computer
	filenameComputer := file.NewNameService()

	var uploaderService uploaderservice.Uploader
	if app.cfg.UploadingStrategy == uploader.MultipartFormUploadingType {
		// used parsing of full form into RAM
		uploaderService =
			uploader.NewNativeUploader(
				logger,
				filesystemStorage,
				filenameComputer,
				app.cfg.ResourceFormFilename,
				app.cfg.InMemoryFileSizeThreshold,
			)
	} else if app.cfg.UploadingStrategy == uploader.MultipartPartUploadingType {
		// used partial reading from multipart.Part
		uploaderService =
			uploader.NewPartsUploader(
				logger,
				filesystemStorage,
				filenameComputer,
			)
	}
	return filesystemStorage, filenameComputer, uploaderService
}

func (app *ResourcesApp) InitRequestResponseServices(
	logger loggerservice.Logger,
) (
	extractor.RequestParams,
	response.Responder,
) {
	req := request.NewParametersExtractor()
	resp := response.NewResponseService(logger)
	return req, resp
}

func (app *ResourcesApp) InitAuthedRestApiControllers(
	cacheService cacheservice.Cacher,
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
			authService,
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
	}
}

func (app *ResourcesApp) InitUnauthedRestApiControllers(
	loggerService loggerservice.Logger,
	responseService response.Responder,
	// user. deps.
	userBuilder builder.User,
	userService userservice.CRUD,
	// auth. deps.
	authBuilder builder.Auth,
	authService authservice.Authenticator,
) []controller.Controller {
	return []controller.Controller{
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

func (app *ResourcesApp) InitAuthedNativeRenderingControllers(
	loggerService loggerservice.Logger,
	responseService response.Responder,
) []controller.Controller {
	return []controller.Controller{
		render.NewIndexController(loggerService, responseService),
	}
}

func (app *ResourcesApp) InitUnauthedNativeRenderingControllers(
	loggerService loggerservice.Logger,
	responseService response.Responder,
) []controller.Controller {
	return []controller.Controller{
		render.NewLoginController(loggerService, responseService),
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
