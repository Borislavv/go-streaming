package di

import (
	"context"
	"github.com/Borislavv/video-streaming/internal/app"
	"github.com/Borislavv/video-streaming/internal/domain/logger"
	"github.com/Borislavv/video-streaming/internal/domain/repository"
	"github.com/Borislavv/video-streaming/internal/domain/service/extractor"
	resourceservice "github.com/Borislavv/video-streaming/internal/domain/service/resource"
	"github.com/Borislavv/video-streaming/internal/domain/service/security"
	"github.com/Borislavv/video-streaming/internal/domain/service/tokenizer"
	"github.com/Borislavv/video-streaming/internal/domain/service/uploader"
	userservice "github.com/Borislavv/video-streaming/internal/domain/service/user"
	videoservice "github.com/Borislavv/video-streaming/internal/domain/service/video"
	"github.com/Borislavv/video-streaming/internal/domain/validator"
	"github.com/Borislavv/video-streaming/internal/infrastructure/api/v1/response"
	"github.com/Borislavv/video-streaming/internal/infrastructure/di"
	"github.com/Borislavv/video-streaming/internal/infrastructure/repository/storage/cache"
	"github.com/Borislavv/video-streaming/internal/infrastructure/repository/storage/mongodb"
	"github.com/Borislavv/video-streaming/internal/infrastructure/service/detector"
	"github.com/Borislavv/video-streaming/internal/infrastructure/service/reader"
	"github.com/Borislavv/video-streaming/internal/infrastructure/service/streamer"
	"github.com/Borislavv/video-streaming/internal/infrastructure/service/streamer/action/handler"
	"github.com/Borislavv/video-streaming/internal/infrastructure/service/streamer/action/handler/strategy"
	"github.com/Borislavv/video-streaming/internal/infrastructure/service/streamer/action/listener"
	"github.com/Borislavv/video-streaming/internal/infrastructure/service/streamer/proto"
	"github.com/Borislavv/video-streaming/internal/infrastructure/service/uploader/file"
	"go.mongodb.org/mongo-driver/mongo"
)

type ContainerManager interface {
	di.Container

	// Application
	GetConfig() (*app.Config, error)
	GetCtx() (context.Context, error)
	GetCancelFunc() (context.CancelFunc, error)

	// Mongo database
	GetMongoDatabase() (*mongo.Database, error)

	// Mongo repository
	GetResourceMongoRepository() (*mongodb.ResourceRepository, error)
	GetVideoMongoRepository() (*mongodb.VideoRepository, error)
	GetUserMongoRepository() (*mongodb.UserRepository, error)

	// Cache repository
	GetResourceCacheRepository() (*cache.ResourceRepository, error)
	GetVideoCacheRepository() (*cache.VideoRepository, error)
	GetUserCacheRepository() (*cache.UserRepository, error)

	// Common services
	GetAccessService() (accesstor_interface.Accessor, error)

	// Resource services
	GetResourceBuilder() (builder_interface.Resource, error)
	GetResourceValidator() (validator.Resource, error)
	GetResourceRepository() (repository.Resource, error)
	GetResourceCRUDService() (resourceservice.CRUD, error)

	// BlockedToken services
	GetBlockedTokenRepository() (repository.BlockedToken, error)

	// Video services
	GetVideoBuilder() (builder_interface.Video, error)
	GetVideoValidator() (validator.Video, error)
	GetVideoRepository() (repository.Video, error)
	GetVideoCRUDService() (videoservice.CRUD, error)

	// User services
	GetUserBuilder() (builder_interface.User, error)
	GetUserValidator() (validator.User, error)
	GetUserRepository() (repository.User, error)
	GetUserCRUDService() (userservice.CRUD, error)

	// Auth services
	GetAuthBuilder() (builder_interface.Auth, error)
	GetAuthValidator() (validator.Auth, error)
	GetAuthService() (authenticator_interface.Authenticator, error)

	// Infrastructure
	GetLoggerService() (logger.Logger, error)
	GetCacheService() (cacher_interface.Cacher, error)
	GetRequestParametersExtractorService() (extractor.RequestParams, error)
	GetResponderService() (response.Responder, error)
	GetPasswordHasherService() (security.PasswordHasher, error)
	GetTokenizerService() (tokenizer.Tokenizer, error)

	// File
	GetFileStorageService() (file.Storage, error)
	GetFileNameComputerService() (file.NameComputer, error)
	GetFileUploaderService() (uploader.Uploader, error)
	GetFileReaderService() (reader.FileReader, error)

	// WebSocket
	GetWebSocketCommunicatorService() (proto.Communicator, error)
	GetWebSocketListener() (listener.ActionsListener, error)
	GetWebSocketHandler() (handler.ActionsHandler, error)
	GetWebSocketHandlerStrategies() ([]strategy.ActionStrategy, error)

	// Codesc
	GetCodecsDetectorService() (detector.Codecs, error)

	// Streaming
	GetStreamingService() (streamer.Streamer, error)
}
