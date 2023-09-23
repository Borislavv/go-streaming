package stream

import (
	"context"
	"github.com/Borislavv/video-streaming/internal/infrastructure/logger"
	"github.com/Borislavv/video-streaming/internal/infrastructure/repository/mongodb"
	"github.com/Borislavv/video-streaming/internal/infrastructure/server/socket"
	"github.com/Borislavv/video-streaming/internal/infrastructure/service/reader"
	"github.com/Borislavv/video-streaming/internal/infrastructure/service/streamer"
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

type StreamingApp struct {
	cfg config
}

func NewStreamingApp() *StreamingApp {
	return &StreamingApp{cfg: config{}}
}

// Run is method which running the streaming part of app
func (app *StreamingApp) Run(mWg *sync.WaitGroup) {
	defer mWg.Done()
	wg := &sync.WaitGroup{}
	ctx, cancel := context.WithCancel(context.Background())

	// init. logger and close func.
	loggerService, cls := logger.NewStdOutLogger(ctx, 10, 10)
	defer func() {
		cancel()
		wg.Wait()
		cls()
	}()

	if err := env.Parse(&app.cfg); err != nil {
		loggerService.Critical(err)
		return
	}

	// init. mongodb client
	mongoClient, err := mongo.Connect(ctx, options.Client().ApplyURI(app.cfg.MongoUri))
	if err != nil {
		loggerService.Critical(err)
		return
	}
	defer func() { _ = mongoClient.Disconnect(ctx) }()

	// ping mongodb
	if err = mongoClient.Ping(ctx, readpref.Primary()); err != nil {
		loggerService.Critical(err)
		return
	}

	// connect to target mongodb database
	db := mongoClient.Database(app.cfg.MongoDb)

	// video repository
	videoRepository := mongodb.NewVideoRepository(db, loggerService, time.Minute)

	// resource reader service
	readerService := reader.NewReaderService(loggerService)

	// custom websocket communication protocol
	wsCommunicator := streamer.NewWebSocketCommunicator(loggerService)

	// resource codecs determiner
	resourceCodecsDetector := streamer.NewResourceCodecInfo(ctx, loggerService)

	// websocket actions listener
	actionsListener := streamer.NewWebSocketActionsListener(loggerService, wsCommunicator)

	// websocket actions handler
	actionsHandler := streamer.NewWebSocketActionsHandler(
		ctx, loggerService, readerService, videoRepository, wsCommunicator, resourceCodecsDetector,
	)

	// resource streaming service
	streamingService := streamer.NewStreamingService(loggerService, actionsListener, actionsHandler)

	wg.Add(1)
	go socket.NewWebSocketServer( // websocket server
		app.cfg.Host, app.cfg.Port, app.cfg.Transport, streamingService, loggerService,
	).Listen(ctx, wg)

	<-app.shutdown()
}

func (app *StreamingApp) shutdown() chan os.Signal {
	stopCh := make(chan os.Signal, 1)
	signal.Notify(stopCh, os.Interrupt, syscall.SIGTERM)
	return stopCh
}
