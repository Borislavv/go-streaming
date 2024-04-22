package resource

import (
	"context"
	"github.com/Borislavv/video-streaming/internal/app"
	"github.com/Borislavv/video-streaming/internal/domain/builder"
	builderinterface "github.com/Borislavv/video-streaming/internal/domain/builder/interface"
	loggerservice "github.com/Borislavv/video-streaming/internal/domain/logger/interface"
	repositoryinterface "github.com/Borislavv/video-streaming/internal/domain/repository/interface"
	"github.com/Borislavv/video-streaming/internal/domain/service/accessor"
	accessorinterface "github.com/Borislavv/video-streaming/internal/domain/service/accessor/interface"
	authservice "github.com/Borislavv/video-streaming/internal/domain/service/authenticator"
	authenticatorinterface "github.com/Borislavv/video-streaming/internal/domain/service/authenticator/interface"
	cacheservice "github.com/Borislavv/video-streaming/internal/domain/service/cacher/interface"
	"github.com/Borislavv/video-streaming/internal/domain/service/di/interface"
	extractorinterface "github.com/Borislavv/video-streaming/internal/domain/service/extractor/interface"
	resourceservice "github.com/Borislavv/video-streaming/internal/domain/service/resource"
	resourceinterface "github.com/Borislavv/video-streaming/internal/domain/service/resource/interface"
	securityservice "github.com/Borislavv/video-streaming/internal/domain/service/security/interface"
	storagerinterface "github.com/Borislavv/video-streaming/internal/domain/service/storager/interface"
	tokenizerinterface "github.com/Borislavv/video-streaming/internal/domain/service/tokenizer/interface"
	uploaderservice "github.com/Borislavv/video-streaming/internal/domain/service/uploader/interface"
	userservice "github.com/Borislavv/video-streaming/internal/domain/service/user"
	userinterface "github.com/Borislavv/video-streaming/internal/domain/service/user/interface"
	videoservice "github.com/Borislavv/video-streaming/internal/domain/service/video"
	videointerface "github.com/Borislavv/video-streaming/internal/domain/service/video/interface"
	"github.com/Borislavv/video-streaming/internal/domain/validator"
	validatorinterface "github.com/Borislavv/video-streaming/internal/domain/validator/interface"
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
	responseinterface "github.com/Borislavv/video-streaming/internal/infrastructure/api/v1/response/interface"
	"github.com/Borislavv/video-streaming/internal/infrastructure/repository/storage/cache"
	"github.com/Borislavv/video-streaming/internal/infrastructure/repository/storage/mongodb"
	mongodbinterface "github.com/Borislavv/video-streaming/internal/infrastructure/repository/storage/mongodb/interface"
	"github.com/Borislavv/video-streaming/internal/infrastructure/server/http"
	"github.com/Borislavv/video-streaming/internal/infrastructure/service/cacher"
	"github.com/Borislavv/video-streaming/internal/infrastructure/service/logger"
	"github.com/Borislavv/video-streaming/internal/infrastructure/service/security"
	"github.com/Borislavv/video-streaming/internal/infrastructure/service/tokenizer"
	"github.com/Borislavv/video-streaming/internal/infrastructure/service/uploader"
	"github.com/Borislavv/video-streaming/internal/infrastructure/service/uploader/file"
	fileinterface "github.com/Borislavv/video-streaming/internal/infrastructure/service/uploader/file/interface"
	"github.com/caarlos0/env/v9"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"os"
	"os/signal"
	"reflect"
	"sync"
	"syscall"
	"time"
)

type ResourcesApp struct {
	cfg *app.Config
	di  diinterface.ServiceContainer
}

func NewResourcesApp(di diinterface.ServiceContainer) *ResourcesApp {
	return &ResourcesApp{cfg: &app.Config{}, di: di}
}

// Run is method which running the REST API part of app
func (app *ResourcesApp) Run(mWg *sync.WaitGroup) {
	defer mWg.Done()
	// ctx, cancelFunc
	cancel := app.InitAppCtx()

	// logger
	loggerService, loggerCancelFunc, err := app.InitLoggerService()
	if err != nil {
		panic(err)
	}
	defer loggerCancelFunc()

	wg := &sync.WaitGroup{}
	defer func() {
		cancel()
		wg.Wait()
	}()

	// config
	if err = app.InitConfig(); err != nil {
		loggerService.Critical(err)
		return
	}

	// mongo database
	databaseCancelFunc, err := app.InitMongoDatabase()
	if err != nil {
		loggerService.Critical(err)
		return
	}
	defer databaseCancelFunc()

	// cache and dependencies
	if err = app.InitCacheService(); err != nil {
		loggerService.Critical(err)
		return
	}

	// request-response and dependencies
	if err = app.InitRequestResponseServices(); err != nil {
		loggerService.Critical(err)
		return
	}

	// access service
	if err = app.InitAccessService(); err != nil {
		loggerService.Critical(err)
		return
	}

	// file uploader and dependencies
	if err = app.InitUploaderServices(); err != nil {
		loggerService.Critical(err)
		return
	}

	// resource services
	if err = app.InitResourceServices(); err != nil {
		loggerService.Critical(err)
		return
	}

	// video services
	if err = app.InitVideoServices(); err != nil {
		loggerService.Critical(err)
		return
	}

	// password services
	if err = app.InitPasswordService(); err != nil {
		loggerService.Critical(err)
		return
	}

	// user services
	if err = app.InitUserServices(); err != nil {
		loggerService.Critical(err)
		return
	}

	// token services
	if err = app.InitTokenServices(); err != nil {
		loggerService.Critical(err)
		return
	}

	// auth services
	if err = app.InitAuthServices(); err != nil {
		loggerService.Critical(err)
		return
	}

	// HTTP server
	if err = app.InitHttpServer(wg); err != nil {
		loggerService.Critical(err)
		return
	}

	<-app.shutdown()
}

func (app *ResourcesApp) shutdown() chan os.Signal {
	stopCh := make(chan os.Signal, 1)
	signal.Notify(stopCh, os.Interrupt, syscall.SIGTERM)
	return stopCh
}

func (app *ResourcesApp) InitAppCtx() context.CancelFunc {
	ctx, cancel := context.WithCancel(context.Background())

	app.di.
		Set(ctx, reflect.TypeOf((*context.Context)(nil))).
		Set(cancel, reflect.TypeOf((*context.CancelFunc)(nil)))

	return cancel
}

func (app *ResourcesApp) InitLoggerService() (
	loggerService loggerservice.Logger,
	deferFunc func(),
	err error,
) {
	ctx, err := app.di.GetCtx()
	if err != nil {
		return nil, nil, err
	}

	loggerService, cls := logger.NewStdOut(ctx, app.cfg.LoggerErrorsBufferCap, app.cfg.LoggerRequestsBufferCap)

	app.di.
		Set(loggerService, reflect.TypeOf((*loggerservice.Logger)(nil))).
		Set(loggerService, nil)

	return loggerService, cls, nil
}

func (app *ResourcesApp) InitConfig() error {
	loggerService, err := app.di.GetLoggerService()
	if err != nil {
		return err
	}

	if err = env.Parse(app.cfg); err != nil {
		return loggerService.LogPropagate(err)
	}

	app.di.
		Set(app.cfg, nil)

	return nil
}

func (app *ResourcesApp) InitMongoDatabase() (deferFunc func(), err error) {
	loggerService, err := app.di.GetLoggerService()
	if err != nil {
		return nil, err
	}

	ctx, err := app.di.GetCtx()
	if err != nil {
		return nil, loggerService.LogPropagate(err)
	}

	c, err := mongo.Connect(ctx, options.Client().ApplyURI(app.cfg.MongoUri))
	if err != nil {
		return nil, loggerService.CriticalPropagate(err)
	}

	deferFunc = func() {
		_ = c.Disconnect(ctx)
	}

	if err = c.Ping(ctx, readpref.Primary()); err != nil {
		return deferFunc, loggerService.CriticalPropagate(err)
	}

	d := c.Database(app.cfg.MongoDb)

	app.di.
		Set(c, nil).
		Set(d, nil)

	return deferFunc, nil
}

func (app *ResourcesApp) InitPasswordService() error {
	loggerService, err := app.di.GetLoggerService()
	if err != nil {
		return err
	}

	cfg, err := app.di.GetConfig()
	if err != nil {
		return loggerService.LogPropagate(err)
	}

	p, err := security.NewPasswordHasher(app.di, cfg.PasswordHashCost)
	if err != nil {
		return loggerService.LogPropagate(err)
	}

	app.di.
		Set(p, reflect.TypeOf((*securityservice.PasswordHasher)(nil))).
		Set(p, nil)

	return nil
}

func (app *ResourcesApp) InitCacheService() error {
	loggerService, err := app.di.GetLoggerService()
	if err != nil {
		return err
	}

	ctx, err := app.di.GetCtx()
	if err != nil {
		return loggerService.LogPropagate(err)
	}

	c := cacher.NewCache(
		cacher.NewMapCacheStorage(ctx),
		cacher.NewCacheDisplacer(ctx, time.Second*1),
	)

	app.di.
		Set(c, reflect.TypeOf((*cacheservice.Cacher)(nil))).
		Set(c, nil)

	return nil
}

func (app *ResourcesApp) InitVideoServices() error {
	loggerService, err := app.di.GetLoggerService()
	if err != nil {
		return err
	}

	r, err := mongodb.NewVideoRepository(app.di)
	if err != nil {
		return loggerService.LogPropagate(err)
	}
	app.di.
		Set(r, reflect.TypeOf((*mongodbinterface.Video)(nil))).
		Set(r, nil)

	c, err := cache.NewVideoRepository(app.di)
	if err != nil {
		return loggerService.LogPropagate(err)
	}
	app.di.
		Set(c, reflect.TypeOf((*repositoryinterface.Video)(nil))).
		Set(c, nil)

	v, err := validator.NewVideoValidator(app.di)
	if err != nil {
		return loggerService.LogPropagate(err)
	}
	app.di.
		Set(v, reflect.TypeOf((*validatorinterface.Video)(nil))).
		Set(v, nil)

	b, err := builder.NewVideoBuilder(app.di)
	if err != nil {
		return loggerService.LogPropagate(err)
	}
	app.di.
		Set(b, reflect.TypeOf((*builderinterface.Video)(nil))).
		Set(b, nil)

	s, err := videoservice.NewCRUDService(app.di)
	if err != nil {
		return loggerService.LogPropagate(err)
	}
	app.di.
		Set(s, reflect.TypeOf((*videointerface.CRUD)(nil))).
		Set(s, nil)

	return nil
}

func (app *ResourcesApp) InitResourceServices() error {
	loggerService, err := app.di.GetLoggerService()
	if err != nil {
		return err
	}

	r, err := mongodb.NewResourceRepository(app.di)
	if err != nil {
		return loggerService.LogPropagate(err)
	}
	app.di.
		Set(r, reflect.TypeOf((*mongodbinterface.Resource)(nil))).
		Set(r, nil)

	c, err := cache.NewResourceRepository(app.di)
	if err != nil {
		return loggerService.LogPropagate(err)
	}
	app.di.
		Set(c, reflect.TypeOf((*repositoryinterface.Resource)(nil))).
		Set(c, nil)

	v, err := validator.NewResourceValidator(app.di)
	if err != nil {
		return loggerService.LogPropagate(err)
	}
	app.di.
		Set(v, reflect.TypeOf((*validatorinterface.Resource)(nil))).
		Set(v, nil)

	b, err := builder.NewResourceBuilder(app.di)
	if err != nil {
		return loggerService.LogPropagate(err)
	}
	app.di.
		Set(b, reflect.TypeOf((*builderinterface.Resource)(nil))).
		Set(b, nil)

	s, err := resourceservice.NewResourceService(app.di)
	if err != nil {
		return loggerService.LogPropagate(err)
	}
	app.di.
		Set(s, reflect.TypeOf((*resourceinterface.CRUD)(nil))).
		Set(s, nil)

	return nil
}

func (app *ResourcesApp) InitUserServices() error {
	loggerService, err := app.di.GetLoggerService()
	if err != nil {
		return err
	}

	r, err := mongodb.NewUserRepository(app.di)
	if err != nil {
		return loggerService.LogPropagate(err)
	}
	app.di.
		Set(r, reflect.TypeOf((*mongodbinterface.User)(nil))).
		Set(r, nil)

	c, err := cache.NewUserRepository(app.di)
	if err != nil {
		return loggerService.LogPropagate(err)
	}
	app.di.
		Set(c, reflect.TypeOf((*repositoryinterface.User)(nil))).
		Set(c, nil)

	b, err := builder.NewUserBuilder(app.di)
	if err != nil {
		return loggerService.LogPropagate(err)
	}
	app.di.
		Set(b, reflect.TypeOf((*builderinterface.User)(nil))).
		Set(b, nil)

	v, err := validator.NewUserValidator(app.di)
	if err != nil {
		return loggerService.LogPropagate(err)
	}
	app.di.
		Set(v, reflect.TypeOf((*validatorinterface.User)(nil))).
		Set(v, nil)

	s, err := userservice.NewCRUDService(app.di)
	if err != nil {
		return loggerService.LogPropagate(err)
	}
	app.di.
		Set(s, reflect.TypeOf((*userinterface.CRUD)(nil))).
		Set(s, nil)

	return nil
}

func (app *ResourcesApp) InitAuthServices() error {
	loggerService, err := app.di.GetLoggerService()
	if err != nil {
		return err
	}

	b, err := builder.NewAuthBuilder(app.di)
	if err != nil {
		return loggerService.LogPropagate(err)
	}
	app.di.
		Set(b, reflect.TypeOf((*builderinterface.Auth)(nil))).
		Set(b, nil)

	v, err := validator.NewAuthValidator(app.di)
	if err != nil {
		return loggerService.LogPropagate(err)
	}
	app.di.
		Set(v, reflect.TypeOf((*validatorinterface.Auth)(nil))).
		Set(v, nil)

	s, err := authservice.NewAuthService(app.di)
	if err != nil {
		return loggerService.LogPropagate(err)
	}
	app.di.
		Set(s, reflect.TypeOf((*authenticatorinterface.Authenticator)(nil))).
		Set(s, nil)

	return nil
}

func (app *ResourcesApp) InitTokenServices() error {
	loggerService, err := app.di.GetLoggerService()
	if err != nil {
		return err
	}

	r, err := mongodb.NewBlockedTokenRepository(app.di)
	if err != nil {
		return loggerService.LogPropagate(err)
	}
	app.di.
		Set(r, reflect.TypeOf((*repositoryinterface.BlockedToken)(nil))).
		Set(r, reflect.TypeOf((*mongodbinterface.BlockedToken)(nil))).
		Set(r, nil)

	s, err := tokenizer.NewJwtService(app.di)
	if err != nil {
		return loggerService.LogPropagate(err)
	}
	app.di.
		Set(s, reflect.TypeOf((*tokenizerinterface.Tokenizer)(nil))).
		Set(s, nil)

	return nil
}

func (app *ResourcesApp) InitUploaderServices() error {
	loggerService, err := app.di.GetLoggerService()
	if err != nil {
		return err
	}

	// filesystem storage
	filesystemStorage, err := file.NewFilesystemStorageService(app.di)
	if err != nil {
		return loggerService.LogPropagate(err)
	}

	// filename computer
	filenameComputer := file.NewNameComputerService()

	app.di.
		Set(filesystemStorage, reflect.TypeOf((*fileinterface.Storage)(nil))).
		Set(filesystemStorage, reflect.TypeOf((*storagerinterface.Storage)(nil))).
		Set(filenameComputer, reflect.TypeOf((*fileinterface.NameComputer)(nil))).
		Set(filesystemStorage, nil).
		Set(filenameComputer, nil)

	if app.cfg.ResourceUploadingStrategy == uploader.MultipartFormUploadingType {
		// used parsing of full form into RAM
		service, uerr := uploader.NewNativeUploader(app.di)
		if uerr != nil {
			return loggerService.LogPropagate(uerr)
		}

		app.di.
			Set(service, reflect.TypeOf((*uploaderservice.Uploader)(nil))).
			Set(service, nil)
	} else if app.cfg.ResourceUploadingStrategy == uploader.MultipartPartUploadingType {
		// used partial reading from multipart.Part
		service, uerr := uploader.NewPartsUploader(app.di)
		if uerr != nil {
			return loggerService.LogPropagate(uerr)
		}

		app.di.
			Set(service, reflect.TypeOf((*uploaderservice.Uploader)(nil))).
			Set(service, nil)
	}

	return nil
}

func (app *ResourcesApp) InitAccessService() error {
	loggerService, err := app.di.GetLoggerService()
	if err != nil {
		return err
	}

	a, err := accessor.NewAccessService(app.di)
	if err != nil {
		return loggerService.LogPropagate(err)
	}

	app.di.
		Set(a, reflect.TypeOf((*accessorinterface.Accessor)(nil))).
		Set(a, nil)

	return nil
}

func (app *ResourcesApp) InitRequestResponseServices() error {
	loggerService, err := app.di.GetLoggerService()
	if err != nil {
		return err
	}

	req := request.NewParametersExtractor()
	resp, err := response.NewResponseService(app.di)
	if err != nil {
		return loggerService.LogPropagate(err)
	}

	app.di.
		Set(req, reflect.TypeOf((*extractorinterface.RequestParams)(nil))).
		Set(resp, reflect.TypeOf((*responseinterface.Responder)(nil))).
		Set(req, nil).
		Set(resp, nil)

	return nil
}

func (app *ResourcesApp) InitAuthedRestApiControllers() ([]controller.Controller, error) {
	loggerService, err := app.di.GetLoggerService()
	if err != nil {
		return nil, err
	}

	// resource
	resourceUploadController, err := resource.NewUploadController(app.di)
	if err != nil {
		return nil, loggerService.LogPropagate(err)
	}

	// user
	userGetController, err := user.NewGetController(app.di)
	if err != nil {
		return nil, loggerService.LogPropagate(err)
	}
	userUpdateController, err := user.NewUpdateUserController(app.di)
	if err != nil {
		return nil, loggerService.LogPropagate(err)
	}
	userDeleteController, err := user.NewDeleteController(app.di)
	if err != nil {
		return nil, loggerService.LogPropagate(err)
	}

	// video
	videoCreateController, err := video.NewCreateController(app.di)
	if err != nil {
		return nil, loggerService.LogPropagate(err)
	}
	videoUpdatedController, err := video.NewUpdateController(app.di)
	if err != nil {
		return nil, loggerService.LogPropagate(err)
	}
	videoGetController, err := video.NewGetController(app.di)
	if err != nil {
		return nil, loggerService.LogPropagate(err)
	}
	videoListController, err := video.NewListController(app.di)
	if err != nil {
		return nil, loggerService.LogPropagate(err)
	}
	videoDeleteController, err := video.NewDeleteController(app.di)
	if err != nil {
		return nil, loggerService.LogPropagate(err)
	}

	return []controller.Controller{
		// resource
		resourceUploadController,
		// video
		videoCreateController,
		videoUpdatedController,
		videoGetController,
		videoListController,
		videoDeleteController,
		// audio
		audio.NewCreateController(),
		audio.NewDeleteController(),
		audio.NewGetController(),
		audio.NewListController(),
		audio.NewUpdateController(),
		// user
		userUpdateController,
		userGetController,
		userDeleteController,
	}, nil
}

func (app *ResourcesApp) InitUnauthedRestApiControllers() ([]controller.Controller, error) {
	loggerService, err := app.di.GetLoggerService()
	if err != nil {
		return nil, err
	}

	authorizationController, err := auth.NewAuthorizationController(app.di)
	if err != nil {
		return nil, loggerService.LogPropagate(err)
	}

	registrationController, err := auth.NewRegistrationController(app.di)
	if err != nil {
		return nil, loggerService.LogPropagate(err)
	}

	return []controller.Controller{
		authorizationController,
		registrationController,
	}, nil
}

func (app *ResourcesApp) InitAuthedNativeRenderingControllers() ([]controller.Controller, error) {
	loggerService, err := app.di.GetLoggerService()
	if err != nil {
		return nil, err
	}

	indexController, err := render.NewIndexController(app.di)
	if err != nil {
		return nil, loggerService.LogPropagate(err)
	}

	return []controller.Controller{
		indexController,
	}, nil
}

func (app *ResourcesApp) InitUnauthedNativeRenderingControllers() ([]controller.Controller, error) {
	loggerService, err := app.di.GetLoggerService()
	if err != nil {
		return nil, err
	}

	loginController, err := render.NewAuthController(app.di)
	if err != nil {
		return nil, loggerService.LogPropagate(err)
	}

	return []controller.Controller{
		loginController,
	}, nil
}

func (app *ResourcesApp) InitStaticServingControllers() ([]controller.Controller, error) {
	loggerService, err := app.di.GetLoggerService()
	if err != nil {
		return nil, err
	}

	staticFilesController, err := static.NewFilesController(app.di)
	if err != nil {
		return nil, loggerService.LogPropagate(err)
	}

	return []controller.Controller{
		staticFilesController,
	}, nil
}

func (app *ResourcesApp) InitHttpServer(wg *sync.WaitGroup) error {
	loggerService, err := app.di.GetLoggerService()
	if err != nil {
		return err
	}

	ctx, err := app.di.GetCtx()
	if err != nil {
		return loggerService.LogPropagate(err)
	}

	// RestAPI
	authedRestAPIControllers, err := app.InitAuthedRestApiControllers()
	if err != nil {
		return loggerService.LogPropagate(err)
	}
	unauthedRestAPIController, err := app.InitUnauthedRestApiControllers()
	if err != nil {
		return loggerService.LogPropagate(err)
	}

	// HTML rendering
	authedNativeControllers, err := app.InitAuthedNativeRenderingControllers()
	if err != nil {
		return loggerService.LogPropagate(err)
	}
	unauthedNativeControllers, err := app.InitUnauthedNativeRenderingControllers()
	if err != nil {
		return loggerService.LogPropagate(err)
	}

	// Static files
	staticFilesControllers, err := app.InitStaticServingControllers()
	if err != nil {
		return loggerService.LogPropagate(err)
	}

	server, err := http.NewHttpServer(
		app.di,
		authedRestAPIControllers,
		unauthedRestAPIController,
		authedNativeControllers,
		unauthedNativeControllers,
		staticFilesControllers,
	)
	if err != nil {
		return loggerService.LogPropagate(err)
	}

	wg.Add(1)
	go server.Listen(ctx, wg)

	return nil
}
