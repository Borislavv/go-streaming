package stream

import (
	"context"
	"github.com/Borislavv/video-streaming/internal/app"
	loggerservice "github.com/Borislavv/video-streaming/internal/domain/logger/interface"
	repositoryinterface "github.com/Borislavv/video-streaming/internal/domain/repository/interface"
	cacheservice "github.com/Borislavv/video-streaming/internal/domain/service/cacher/interface"
	"github.com/Borislavv/video-streaming/internal/domain/service/di/interface"
	tokenizerinterface "github.com/Borislavv/video-streaming/internal/domain/service/tokenizer/interface"
	"github.com/Borislavv/video-streaming/internal/infrastructure/repository/storage/cache"
	"github.com/Borislavv/video-streaming/internal/infrastructure/repository/storage/mongodb"
	mongodbinterface "github.com/Borislavv/video-streaming/internal/infrastructure/repository/storage/mongodb/interface"
	server "github.com/Borislavv/video-streaming/internal/infrastructure/server/ws"
	"github.com/Borislavv/video-streaming/internal/infrastructure/service/cacher"
	"github.com/Borislavv/video-streaming/internal/infrastructure/service/detector"
	detectorinterface "github.com/Borislavv/video-streaming/internal/infrastructure/service/detector/interface"
	"github.com/Borislavv/video-streaming/internal/infrastructure/service/logger"
	"github.com/Borislavv/video-streaming/internal/infrastructure/service/reader"
	readerinterface "github.com/Borislavv/video-streaming/internal/infrastructure/service/reader/interface"
	"github.com/Borislavv/video-streaming/internal/infrastructure/service/streamer"
	"github.com/Borislavv/video-streaming/internal/infrastructure/service/streamer/action/handler"
	handlerinterface "github.com/Borislavv/video-streaming/internal/infrastructure/service/streamer/action/handler/interface"
	"github.com/Borislavv/video-streaming/internal/infrastructure/service/streamer/action/handler/strategy"
	strategyinterface "github.com/Borislavv/video-streaming/internal/infrastructure/service/streamer/action/handler/strategy/interface"
	"github.com/Borislavv/video-streaming/internal/infrastructure/service/streamer/action/listener"
	listenerinterface "github.com/Borislavv/video-streaming/internal/infrastructure/service/streamer/action/listener/interface"
	streamerinterface "github.com/Borislavv/video-streaming/internal/infrastructure/service/streamer/interface"
	protointerface "github.com/Borislavv/video-streaming/internal/infrastructure/service/streamer/proto/interface"
	"github.com/Borislavv/video-streaming/internal/infrastructure/service/streamer/proto/ws"
	"github.com/Borislavv/video-streaming/internal/infrastructure/service/tokenizer"
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

type StreamingApp struct {
	cfg *app.Config
	di  diinterface.ServiceContainer
}

func NewStreamingApp(di diinterface.ServiceContainer) *StreamingApp {
	return &StreamingApp{cfg: &app.Config{}, di: di}
}

// Run is method which running the streaming part of app
func (app *StreamingApp) Run(mWg *sync.WaitGroup) {
	defer mWg.Done()

	cancel := app.InitAppCtx()

	loggerService, loggerCancelFunc, err := app.InitLoggerService()
	if err != nil {
		panic(err)
	}
	defer loggerCancelFunc()

	wg := &sync.WaitGroup{}
	defer func() {
		cancel()
		wg.Wait()
		time.Sleep(time.Second)
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

	// video services
	if err = app.InitVideoServices(); err != nil {
		loggerService.Critical(err)
		return
	}

	// file reader service
	if err = app.InitFileReaderService(); err != nil {
		loggerService.Critical(err)
		return
	}

	// websocket communication service
	if err = app.InitWebSocketCommunicator(); err != nil {
		loggerService.Critical(err)
		return
	}

	// resource codecs detector service
	if err = app.InitCodecsInfoService(); err != nil {
		loggerService.Critical(err)
		return
	}

	// token services
	if err = app.InitTokenServices(); err != nil {
		loggerService.Critical(err)
		return
	}

	// websocket actions listener
	if err = app.InitWebSocketListener(); err != nil {
		loggerService.Critical(err)
		return
	}

	// websocket actions handler
	if err = app.InitWebSocketHandler(); err != nil {
		loggerService.Critical(err)
		return
	}

	// resource streaming service
	if err = app.InitStreamingService(); err != nil {
		loggerService.Critical(err)
		return
	}

	// WebSocket server
	if err = app.InitWebSocketServer(wg); err != nil {
		loggerService.Critical(err)
		return
	}

	<-app.shutdown()
}

func (app *StreamingApp) shutdown() chan os.Signal {
	stopCh := make(chan os.Signal, 1)
	signal.Notify(stopCh, syscall.SIGINT, syscall.SIGTERM, syscall.SIGSTOP)
	return stopCh
}

func (app *StreamingApp) InitAppCtx() context.CancelFunc {
	ctx, cancel := context.WithCancel(context.Background())

	app.di.
		Set(ctx, reflect.TypeOf((*context.Context)(nil))).
		Set(cancel, reflect.TypeOf((*context.CancelFunc)(nil)))

	return cancel
}

func (app *StreamingApp) InitLoggerService() (
	loggerService loggerservice.Logger,
	deferFunc func(),
	err error,
) {
	ctx, err := app.di.GetCtx()
	if err != nil {
		return nil, nil, err
	}

	loggerService, cls := logger.NewStdErr(ctx, app.cfg.LoggerErrorsBufferCap, app.cfg.LoggerRequestsBufferCap)

	app.di.
		Set(loggerService, reflect.TypeOf((*loggerservice.Logger)(nil))).
		Set(loggerService, nil)

	return loggerService, cls, nil
}

func (app *StreamingApp) InitConfig() error {
	if err := env.Parse(app.cfg); err != nil {
		return err
	}

	app.di.
		Set(app.cfg, nil)

	return nil
}

func (app *StreamingApp) InitMongoDatabase() (deferFunc func(), err error) {
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

func (app *StreamingApp) InitCacheService() error {
	ctx, err := app.di.GetCtx()
	if err != nil {
		return err
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

func (app *StreamingApp) InitVideoServices() error {
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

	return nil
}

func (app *StreamingApp) InitFileReaderService() error {
	loggerService, err := app.di.GetLoggerService()
	if err != nil {
		return err
	}

	r, err := reader.NewFileReaderService(app.di)
	if err != nil {
		return loggerService.LogPropagate(err)
	}
	app.di.
		Set(r, reflect.TypeOf((*readerinterface.FileReader)(nil))).
		Set(r, nil)

	return nil
}

func (app *StreamingApp) InitWebSocketCommunicator() error {
	loggerService, err := app.di.GetLoggerService()
	if err != nil {
		return err
	}

	c, err := ws.NewWebSocketCommunicator(app.di)
	if err != nil {
		return loggerService.LogPropagate(err)
	}
	app.di.
		Set(c, reflect.TypeOf((*protointerface.Communicator)(nil))).
		Set(c, nil)

	return nil
}

func (app *StreamingApp) InitCodecsInfoService() error {
	loggerService, err := app.di.GetLoggerService()
	if err != nil {
		return err
	}

	c, err := detector.NewResourceCodecs(app.di)
	if err != nil {
		return loggerService.LogPropagate(err)
	}
	app.di.
		Set(c, reflect.TypeOf((*detectorinterface.Codecs)(nil))).
		Set(c, nil)

	return nil
}

func (app *StreamingApp) InitWebSocketListener() error {
	loggerService, err := app.di.GetLoggerService()
	if err != nil {
		return err
	}

	l, err := listener.NewWebSocketActionsListener(app.di)
	if err != nil {
		return loggerService.LogPropagate(err)
	}
	app.di.
		Set(l, reflect.TypeOf((*listenerinterface.ActionsListener)(nil))).
		Set(l, nil)

	return nil
}

func (app *StreamingApp) InitWebSocketHandler() error {
	loggerService, err := app.di.GetLoggerService()
	if err != nil {
		return err
	}

	// strategies
	streamByIDStrategy, err := strategy.NewStreamByIDActionStrategy(app.di)
	if err != nil {
		return loggerService.LogPropagate(err)
	}
	app.di.
		Set(streamByIDStrategy, nil).
		Set([]strategyinterface.ActionStrategy{
			streamByIDStrategy,
		}, reflect.TypeOf((*[]strategyinterface.ActionStrategy)(nil)))

	// handler which use strategies
	h, err := handler.NewWebSocketActionsHandler(app.di)
	if err != nil {
		return loggerService.LogPropagate(err)
	}
	app.di.
		Set(h, reflect.TypeOf((*handlerinterface.ActionsHandler)(nil))).
		Set(h, nil)

	return nil
}

func (app *StreamingApp) InitStreamingService() error {
	loggerService, err := app.di.GetLoggerService()
	if err != nil {
		return err
	}

	s, err := streamer.NewStreamingService(app.di)
	if err != nil {
		return loggerService.LogPropagate(err)
	}
	app.di.
		Set(s, reflect.TypeOf((*streamerinterface.Streamer)(nil))).
		Set(s, nil)

	return nil
}

func (app *StreamingApp) InitTokenServices() error {
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

func (app *StreamingApp) InitWebSocketServer(wg *sync.WaitGroup) error {
	loggerService, err := app.di.GetLoggerService()
	if err != nil {
		return err
	}

	ctx, err := app.di.GetCtx()
	if err != nil {
		return loggerService.LogPropagate(err)
	}

	s, err := server.NewWebSocketServer(app.di)
	if err != nil {
		return loggerService.LogPropagate(err)
	}

	wg.Add(1)
	go s.Listen(ctx, wg)

	return nil
}
