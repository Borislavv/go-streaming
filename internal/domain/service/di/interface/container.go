package diinterface

import (
	"context"
	"github.com/Borislavv/video-streaming/internal/app"
	builder_interface "github.com/Borislavv/video-streaming/internal/domain/builder/interface"
	"github.com/Borislavv/video-streaming/internal/domain/logger/interface"
	repository_interface "github.com/Borislavv/video-streaming/internal/domain/repository/interface"
	accessor_interface "github.com/Borislavv/video-streaming/internal/domain/service/accessor/interface"
	authenticator_interface "github.com/Borislavv/video-streaming/internal/domain/service/authenticator/interface"
	cacher_interface "github.com/Borislavv/video-streaming/internal/domain/service/cacher/interface"
	extractor_interface "github.com/Borislavv/video-streaming/internal/domain/service/extractor/interface"
	resourceservice "github.com/Borislavv/video-streaming/internal/domain/service/resource/interface"
	security_interface "github.com/Borislavv/video-streaming/internal/domain/service/security/interface"
	tokenizer_interface "github.com/Borislavv/video-streaming/internal/domain/service/tokenizer/interface"
	uploader_interface "github.com/Borislavv/video-streaming/internal/domain/service/uploader/interface"
	userservice "github.com/Borislavv/video-streaming/internal/domain/service/user/interface"
	videoservice "github.com/Borislavv/video-streaming/internal/domain/service/video/interface"
	validator_interface "github.com/Borislavv/video-streaming/internal/domain/validator/interface"
	response_interface "github.com/Borislavv/video-streaming/internal/infrastructure/api/v1/response/interface"
	diinterface "github.com/Borislavv/video-streaming/internal/infrastructure/di/interface"
	cache_interface "github.com/Borislavv/video-streaming/internal/infrastructure/repository/storage/cache/interface"
	mongodbinterface "github.com/Borislavv/video-streaming/internal/infrastructure/repository/storage/mongodb/interface"
	detector_interface "github.com/Borislavv/video-streaming/internal/infrastructure/service/detector/interface"
	reader_interface "github.com/Borislavv/video-streaming/internal/infrastructure/service/reader/interface"
	handler_interface "github.com/Borislavv/video-streaming/internal/infrastructure/service/streamer/action/handler/interface"
	strategy_interface "github.com/Borislavv/video-streaming/internal/infrastructure/service/streamer/action/handler/strategy/interface"
	listener_interface "github.com/Borislavv/video-streaming/internal/infrastructure/service/streamer/action/listener/interface"
	streamer_interface "github.com/Borislavv/video-streaming/internal/infrastructure/service/streamer/interface"
	proto_interface "github.com/Borislavv/video-streaming/internal/infrastructure/service/streamer/proto/interface"
	file_interface "github.com/Borislavv/video-streaming/internal/infrastructure/service/uploader/file/interface"
	"go.mongodb.org/mongo-driver/mongo"
)

type ContainerManager interface {
	diinterface.Container

	// Application
	GetConfig() (*app.Config, error)
	GetCtx() (context.Context, error)
	GetCancelFunc() (context.CancelFunc, error)

	// Mongo database
	GetMongoDatabase() (*mongo.Database, error)

	// Mongo repository
	GetResourceMongoRepository() (mongodbinterface.Resource, error)
	GetVideoMongoRepository() (mongodbinterface.Video, error)
	GetUserMongoRepository() (mongodbinterface.User, error)
	GetBlockedTokenMongoRepository() (mongodbinterface.BlockedToken, error)

	// Cache repository
	GetResourceCacheRepository() (cache_interface.Resource, error)
	GetVideoCacheRepository() (cache_interface.Video, error)
	GetUserCacheRepository() (cache_interface.User, error)

	// Common services
	GetAccessService() (accessor_interface.Accessor, error)

	// Resource services
	GetResourceBuilder() (builder_interface.Resource, error)
	GetResourceValidator() (validator_interface.Resource, error)
	GetResourceRepository() (repository_interface.Resource, error)
	GetResourceCRUDService() (resourceservice.CRUD, error)

	// BlockedToken services
	GetBlockedTokenRepository() (repository_interface.BlockedToken, error)

	// Video services
	GetVideoBuilder() (builder_interface.Video, error)
	GetVideoValidator() (validator_interface.Video, error)
	GetVideoRepository() (repository_interface.Video, error)
	GetVideoCRUDService() (videoservice.CRUD, error)

	// User services
	GetUserBuilder() (builder_interface.User, error)
	GetUserValidator() (validator_interface.User, error)
	GetUserRepository() (repository_interface.User, error)
	GetUserCRUDService() (userservice.CRUD, error)

	// Auth services
	GetAuthBuilder() (builder_interface.Auth, error)
	GetAuthValidator() (validator_interface.Auth, error)
	GetAuthService() (authenticator_interface.Authenticator, error)

	// Infrastructure
	GetLoggerService() (logger_interface.Logger, error)
	GetCacheService() (cacher_interface.Cacher, error)
	GetRequestParametersExtractorService() (extractor_interface.RequestParams, error)
	GetResponderService() (response_interface.Responder, error)
	GetPasswordHasherService() (security_interface.PasswordHasher, error)
	GetTokenizerService() (tokenizer_interface.Tokenizer, error)

	// File
	GetFileStorageService() (file_interface.Storage, error)
	GetFileNameComputerService() (file_interface.NameComputer, error)
	GetFileUploaderService() (uploader_interface.Uploader, error)
	GetFileReaderService() (reader_interface.FileReader, error)

	// WebSocket
	GetWebSocketCommunicatorService() (proto_interface.Communicator, error)
	GetWebSocketListener() (listener_interface.ActionsListener, error)
	GetWebSocketHandler() (handler_interface.ActionsHandler, error)
	GetWebSocketHandlerStrategies() ([]strategy_interface.ActionStrategy, error)

	// Codesc
	GetCodecsDetectorService() (detector_interface.Codecs, error)

	// Streaming
	GetStreamingService() (streamer_interface.Streamer, error)
}
