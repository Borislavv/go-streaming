package stream

import (
	"context"
	loggerservice "github.com/Borislavv/video-streaming/internal/domain/logger"
	"github.com/Borislavv/video-streaming/internal/domain/repository"
	tokenizerservice "github.com/Borislavv/video-streaming/internal/domain/service/tokenizer"
	"github.com/Borislavv/video-streaming/internal/infrastructure/repository/storage/mongodb"
	server "github.com/Borislavv/video-streaming/internal/infrastructure/server/ws"
	"github.com/Borislavv/video-streaming/internal/infrastructure/service/detector"
	"github.com/Borislavv/video-streaming/internal/infrastructure/service/logger"
	"github.com/Borislavv/video-streaming/internal/infrastructure/service/reader"
	"github.com/Borislavv/video-streaming/internal/infrastructure/service/streamer"
	"github.com/Borislavv/video-streaming/internal/infrastructure/service/streamer/action/handler"
	"github.com/Borislavv/video-streaming/internal/infrastructure/service/streamer/action/handler/strategy"
	"github.com/Borislavv/video-streaming/internal/infrastructure/service/streamer/action/listener"
	"github.com/Borislavv/video-streaming/internal/infrastructure/service/streamer/proto/ws"
	"github.com/Borislavv/video-streaming/internal/infrastructure/service/tokenizer"
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
	loggerService, cls := logger.NewStdOut(ctx, 10, 10)
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
	readerService := reader.NewFileReaderService(ctx, loggerService, app.cfg.ChunkSize)

	// custom websocket communication protocol
	wsCommunicator := ws.NewWebSocketCommunicator(loggerService)

	// resource codecs determiner
	codecsDetector := detector.NewResourceCodecInfo(ctx, loggerService)

	// websocket actions listener
	actionsListener := listener.NewWebSocketActionsListener(loggerService, wsCommunicator)

	// Tokenizer service
	_, tokenService := app.InitTokenServices(ctx, loggerService, db)

	// websocket actions handler
	actionsHandler := handler.NewWebSocketActionsHandler(
		ctx,
		loggerService,
		[]strategy.ActionStrategy{
			strategy.NewStreamByIDActionStrategy(
				ctx, loggerService, videoRepository, readerService,
				codecsDetector, wsCommunicator, tokenService,
			),
		},
	)

	// resource streaming service
	streamingService := streamer.NewStreamingService(loggerService, actionsListener, actionsHandler)

	wg.Add(1)
	go server.NewWebSocketServer( // websocket server
		app.cfg.Host, app.cfg.Port, app.cfg.Transport, streamingService, loggerService,
	).Listen(ctx, wg)

	<-app.shutdown()
}

func (app *StreamingApp) shutdown() chan os.Signal {
	stopCh := make(chan os.Signal, 1)
	signal.Notify(stopCh, os.Interrupt, syscall.SIGTERM)
	return stopCh
}

func (app *StreamingApp) InitTokenServices(
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
