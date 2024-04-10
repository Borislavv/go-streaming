package diinterface

import (
	"context"
	"github.com/Borislavv/video-streaming/internal/app"
	builderinterface "github.com/Borislavv/video-streaming/internal/domain/builder/interface"
	"github.com/Borislavv/video-streaming/internal/domain/logger/interface"
	repositoryinterface "github.com/Borislavv/video-streaming/internal/domain/repository/interface"
	accessorinterface "github.com/Borislavv/video-streaming/internal/domain/service/accessor/interface"
	authenticatorinterface "github.com/Borislavv/video-streaming/internal/domain/service/authenticator/interface"
	cacherinterface "github.com/Borislavv/video-streaming/internal/domain/service/cacher/interface"
	extractorinterface "github.com/Borislavv/video-streaming/internal/domain/service/extractor/interface"
	resourceservice "github.com/Borislavv/video-streaming/internal/domain/service/resource/interface"
	securityinterface "github.com/Borislavv/video-streaming/internal/domain/service/security/interface"
	tokenizerinterface "github.com/Borislavv/video-streaming/internal/domain/service/tokenizer/interface"
	uploaderinterface "github.com/Borislavv/video-streaming/internal/domain/service/uploader/interface"
	userservice "github.com/Borislavv/video-streaming/internal/domain/service/user/interface"
	videoservice "github.com/Borislavv/video-streaming/internal/domain/service/video/interface"
	validatorinterface "github.com/Borislavv/video-streaming/internal/domain/validator/interface"
	responseinterface "github.com/Borislavv/video-streaming/internal/infrastructure/api/v1/response/interface"
	diinterface "github.com/Borislavv/video-streaming/internal/infrastructure/di/interface"
	cacheinterface "github.com/Borislavv/video-streaming/internal/infrastructure/repository/storage/cache/interface"
	mongodbinterface "github.com/Borislavv/video-streaming/internal/infrastructure/repository/storage/mongodb/interface"
	detectorinterface "github.com/Borislavv/video-streaming/internal/infrastructure/service/detector/interface"
	readerinterface "github.com/Borislavv/video-streaming/internal/infrastructure/service/reader/interface"
	handlerinterface "github.com/Borislavv/video-streaming/internal/infrastructure/service/streamer/action/handler/interface"
	strategyinterface "github.com/Borislavv/video-streaming/internal/infrastructure/service/streamer/action/handler/strategy/interface"
	listenerinterface "github.com/Borislavv/video-streaming/internal/infrastructure/service/streamer/action/listener/interface"
	streamerinterface "github.com/Borislavv/video-streaming/internal/infrastructure/service/streamer/interface"
	protointerface "github.com/Borislavv/video-streaming/internal/infrastructure/service/streamer/proto/interface"
	fileinterface "github.com/Borislavv/video-streaming/internal/infrastructure/service/uploader/file/interface"
	"go.mongodb.org/mongo-driver/mongo"
)

type ServiceContainer interface {
	diinterface.Container

	GetConfig() (*app.Config, error)
	GetCtx() (context.Context, error)
	GetCancelFunc() (context.CancelFunc, error)

	GetMongoDatabase() (*mongo.Database, error)

	GetResourceMongoRepository() (mongodbinterface.Resource, error)
	GetVideoMongoRepository() (mongodbinterface.Video, error)
	GetUserMongoRepository() (mongodbinterface.User, error)
	GetBlockedTokenMongoRepository() (mongodbinterface.BlockedToken, error)

	GetResourceCacheRepository() (cacheinterface.Resource, error)
	GetVideoCacheRepository() (cacheinterface.Video, error)
	GetUserCacheRepository() (cacheinterface.User, error)

	GetAccessService() (accessorinterface.Accessor, error)

	GetResourceBuilder() (builderinterface.Resource, error)
	GetResourceValidator() (validatorinterface.Resource, error)
	GetResourceRepository() (repositoryinterface.Resource, error)
	GetResourceCRUDService() (resourceservice.CRUD, error)

	GetBlockedTokenRepository() (repositoryinterface.BlockedToken, error)

	GetVideoBuilder() (builderinterface.Video, error)
	GetVideoValidator() (validatorinterface.Video, error)
	GetVideoRepository() (repositoryinterface.Video, error)
	GetVideoCRUDService() (videoservice.CRUD, error)

	GetUserBuilder() (builderinterface.User, error)
	GetUserValidator() (validatorinterface.User, error)
	GetUserRepository() (repositoryinterface.User, error)
	GetUserCRUDService() (userservice.CRUD, error)

	GetAuthBuilder() (builderinterface.Auth, error)
	GetAuthValidator() (validatorinterface.Auth, error)
	GetAuthService() (authenticatorinterface.Authenticator, error)

	GetLoggerService() (loggerinterface.Logger, error)
	GetCacheService() (cacherinterface.Cacher, error)
	GetRequestParametersExtractorService() (extractorinterface.RequestParams, error)
	GetResponderService() (responseinterface.Responder, error)
	GetPasswordHasherService() (securityinterface.PasswordHasher, error)
	GetTokenizerService() (tokenizerinterface.Tokenizer, error)

	GetFileStorageService() (fileinterface.Storage, error)
	GetFileNameComputerService() (fileinterface.NameComputer, error)
	GetFileUploaderService() (uploaderinterface.Uploader, error)
	GetFileReaderService() (readerinterface.FileReader, error)

	GetWebSocketCommunicatorService() (protointerface.Communicator, error)
	GetWebSocketListener() (listenerinterface.ActionsListener, error)
	GetWebSocketHandler() (handlerinterface.ActionsHandler, error)
	GetWebSocketHandlerStrategies() ([]strategyinterface.ActionStrategy, error)

	GetCodecsDetectorService() (detectorinterface.Codecs, error)

	GetStreamingService() (streamerinterface.Streamer, error)
}
