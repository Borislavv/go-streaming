package di

import (
	"context"
	"github.com/Borislavv/video-streaming/internal/app"
	builderinterface "github.com/Borislavv/video-streaming/internal/domain/builder/interface"
	"github.com/Borislavv/video-streaming/internal/domain/errtype"
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
	"reflect"
)

type ServiceContainer struct {
	diinterface.Container
}

func NewServiceContainerManager() *ServiceContainer {
	return &ServiceContainer{
		Container: NewServiceContainer(),
	}
}

func (s *ServiceContainer) GetConfig() (*app.Config, error) {
	key := (*app.Config)(nil)
	service, err := s.Get(reflect.TypeOf(key))
	if err != nil {
		return nil, errtype.NewServiceWasNotFoundIntoContainerError(reflect.TypeOf(key))
	}
	cfg, ok := service.Interface().(*app.Config)
	if !ok {
		return nil, errtype.NewTypesMismatchedServiceContainerError(reflect.TypeOf(service), reflect.TypeOf(key))
	}
	return cfg, nil
}

func (s *ServiceContainer) GetCtx() (context.Context, error) {
	key := (*context.Context)(nil)
	service, err := s.Get(reflect.TypeOf(key))
	if err != nil {
		return nil, errtype.NewServiceWasNotFoundIntoContainerError(reflect.TypeOf(key))
	}
	ctx, ok := service.Interface().(context.Context)
	if !ok {
		return nil, errtype.NewTypesMismatchedServiceContainerError(reflect.TypeOf(service), reflect.TypeOf(key))
	}
	return ctx, nil
}

func (s *ServiceContainer) GetCancelFunc() (context.CancelFunc, error) {
	key := (*context.CancelFunc)(nil)
	service, err := s.Get(reflect.TypeOf(key))
	if err != nil {
		return nil, errtype.NewServiceWasNotFoundIntoContainerError(reflect.TypeOf(key))
	}
	cancel, ok := service.Interface().(context.CancelFunc)
	if !ok {
		return nil, errtype.NewTypesMismatchedServiceContainerError(reflect.TypeOf(service), reflect.TypeOf(key))
	}
	return cancel, nil
}

func (s *ServiceContainer) GetMongoDatabase() (*mongo.Database, error) {
	key := (*mongo.Database)(nil)
	service, err := s.Get(reflect.TypeOf(key))
	if err != nil {
		return nil, errtype.NewServiceWasNotFoundIntoContainerError(reflect.TypeOf(key))
	}
	database, ok := service.Interface().(*mongo.Database)
	if !ok {
		return nil, errtype.NewTypesMismatchedServiceContainerError(reflect.TypeOf(service), reflect.TypeOf(key))
	}
	return database, nil
}

func (s *ServiceContainer) GetResourceMongoRepository() (mongodbinterface.Resource, error) {
	key := (*mongodbinterface.Resource)(nil)
	service, err := s.Get(reflect.TypeOf(key))
	if err != nil {
		return nil, errtype.NewServiceWasNotFoundIntoContainerError(reflect.TypeOf(key))
	}
	repo, ok := service.Interface().(mongodbinterface.Resource)
	if !ok {
		return nil, errtype.NewTypesMismatchedServiceContainerError(reflect.TypeOf(service), reflect.TypeOf(key))
	}
	return repo, nil
}

func (s *ServiceContainer) GetVideoMongoRepository() (mongodbinterface.Video, error) {
	key := (*mongodbinterface.Video)(nil)
	service, err := s.Get(reflect.TypeOf(key))
	if err != nil {
		return nil, errtype.NewServiceWasNotFoundIntoContainerError(reflect.TypeOf(key))
	}
	repo, ok := service.Interface().(mongodbinterface.Video)
	if !ok {
		return nil, errtype.NewTypesMismatchedServiceContainerError(reflect.TypeOf(service), reflect.TypeOf(key))
	}
	return repo, nil
}

func (s *ServiceContainer) GetUserMongoRepository() (mongodbinterface.User, error) {
	key := (*mongodbinterface.User)(nil)
	service, err := s.Get(reflect.TypeOf(key))
	if err != nil {
		return nil, errtype.NewServiceWasNotFoundIntoContainerError(reflect.TypeOf(key))
	}
	repo, ok := service.Interface().(mongodbinterface.User)
	if !ok {
		return nil, errtype.NewTypesMismatchedServiceContainerError(reflect.TypeOf(service), reflect.TypeOf(key))
	}
	return repo, nil
}

func (s *ServiceContainer) GetBlockedTokenMongoRepository() (mongodbinterface.BlockedToken, error) {
	key := (*mongodbinterface.BlockedToken)(nil)
	service, err := s.Get(reflect.TypeOf(key))
	if err != nil {
		return nil, errtype.NewServiceWasNotFoundIntoContainerError(reflect.TypeOf(key))
	}
	repo, ok := service.Interface().(mongodbinterface.BlockedToken)
	if !ok {
		return nil, errtype.NewTypesMismatchedServiceContainerError(reflect.TypeOf(service), reflect.TypeOf(key))
	}
	return repo, nil
}

func (s *ServiceContainer) GetResourceCacheRepository() (cacheinterface.Resource, error) {
	key := (*cacheinterface.Resource)(nil)
	service, err := s.Get(reflect.TypeOf(key))
	if err != nil {
		return nil, errtype.NewServiceWasNotFoundIntoContainerError(reflect.TypeOf(key))
	}
	repo, ok := service.Interface().(cacheinterface.Resource)
	if !ok {
		return nil, errtype.NewTypesMismatchedServiceContainerError(reflect.TypeOf(service), reflect.TypeOf(key))
	}
	return repo, nil
}

func (s *ServiceContainer) GetVideoCacheRepository() (cacheinterface.Video, error) {
	key := (*cacheinterface.Video)(nil)
	service, err := s.Get(reflect.TypeOf(key))
	if err != nil {
		return nil, errtype.NewServiceWasNotFoundIntoContainerError(reflect.TypeOf(key))
	}
	repo, ok := service.Interface().(cacheinterface.Video)
	if !ok {
		return nil, errtype.NewTypesMismatchedServiceContainerError(reflect.TypeOf(service), reflect.TypeOf(key))
	}
	return repo, nil
}

func (s *ServiceContainer) GetUserCacheRepository() (cacheinterface.User, error) {
	key := (*cacheinterface.User)(nil)
	service, err := s.Get(reflect.TypeOf(key))
	if err != nil {
		return nil, errtype.NewServiceWasNotFoundIntoContainerError(reflect.TypeOf(key))
	}
	repo, ok := service.Interface().(cacheinterface.User)
	if !ok {
		return nil, errtype.NewTypesMismatchedServiceContainerError(reflect.TypeOf(service), reflect.TypeOf(key))
	}
	return repo, nil
}

func (s *ServiceContainer) GetAccessService() (accessorinterface.Accessor, error) {
	key := (*accessorinterface.Accessor)(nil)
	service, err := s.Get(reflect.TypeOf(key))
	if err != nil {
		return nil, errtype.NewServiceWasNotFoundIntoContainerError(reflect.TypeOf(key))
	}
	accessService, ok := service.Interface().(accessorinterface.Accessor)
	if !ok {
		return nil, errtype.NewTypesMismatchedServiceContainerError(reflect.TypeOf(service), reflect.TypeOf(key))
	}
	return accessService, nil
}

func (s *ServiceContainer) GetResourceBuilder() (builderinterface.Resource, error) {
	key := (*builderinterface.Resource)(nil)
	service, err := s.Get(reflect.TypeOf(key))
	if err != nil {
		return nil, errtype.NewServiceWasNotFoundIntoContainerError(reflect.TypeOf(key))
	}
	builderService, ok := service.Interface().(builderinterface.Resource)
	if !ok {
		return nil, errtype.NewTypesMismatchedServiceContainerError(reflect.TypeOf(service), reflect.TypeOf(key))
	}
	return builderService, nil
}

func (s *ServiceContainer) GetResourceValidator() (validatorinterface.Resource, error) {
	key := (*validatorinterface.Resource)(nil)
	service, err := s.Get(reflect.TypeOf(key))
	if err != nil {
		return nil, errtype.NewServiceWasNotFoundIntoContainerError(reflect.TypeOf(key))
	}
	validatorService, ok := service.Interface().(validatorinterface.Resource)
	if !ok {
		return nil, errtype.NewTypesMismatchedServiceContainerError(reflect.TypeOf(service), reflect.TypeOf(key))
	}
	return validatorService, nil
}

func (s *ServiceContainer) GetResourceRepository() (repositoryinterface.Resource, error) {
	key := (*repositoryinterface.Resource)(nil)
	service, err := s.Get(reflect.TypeOf(key))
	if err != nil {
		return nil, errtype.NewServiceWasNotFoundIntoContainerError(reflect.TypeOf(key))
	}
	repo, ok := service.Interface().(repositoryinterface.Resource)
	if !ok {
		return nil, errtype.NewTypesMismatchedServiceContainerError(reflect.TypeOf(service), reflect.TypeOf(key))
	}
	return repo, nil
}

func (s *ServiceContainer) GetResourceCRUDService() (resourceservice.CRUD, error) {
	key := (*resourceservice.CRUD)(nil)
	reflectService, err := s.Get(reflect.TypeOf(key))
	if err != nil {
		return nil, errtype.NewServiceWasNotFoundIntoContainerError(reflect.TypeOf(key))
	}
	service, ok := reflectService.Interface().(resourceservice.CRUD)
	if !ok {
		return nil, errtype.NewTypesMismatchedServiceContainerError(reflect.TypeOf(reflectService), reflect.TypeOf(key))
	}
	return service, nil
}

func (s *ServiceContainer) GetBlockedTokenRepository() (repositoryinterface.BlockedToken, error) {
	key := (*repositoryinterface.BlockedToken)(nil)
	reflectService, err := s.Get(reflect.TypeOf(key))
	if err != nil {
		return nil, errtype.NewServiceWasNotFoundIntoContainerError(reflect.TypeOf(key))
	}
	service, ok := reflectService.Interface().(repositoryinterface.BlockedToken)
	if !ok {
		return nil, errtype.NewTypesMismatchedServiceContainerError(reflect.TypeOf(reflectService), reflect.TypeOf(key))
	}
	return service, nil
}

func (s *ServiceContainer) GetVideoBuilder() (builderinterface.Video, error) {
	key := (*builderinterface.Video)(nil)
	reflectService, err := s.Get(reflect.TypeOf(key))
	if err != nil {
		return nil, errtype.NewServiceWasNotFoundIntoContainerError(reflect.TypeOf(key))
	}
	service, ok := reflectService.Interface().(builderinterface.Video)
	if !ok {
		return nil, errtype.NewTypesMismatchedServiceContainerError(reflect.TypeOf(reflectService), reflect.TypeOf(key))
	}
	return service, nil
}

func (s *ServiceContainer) GetVideoValidator() (validatorinterface.Video, error) {
	key := (*validatorinterface.Video)(nil)
	reflectService, err := s.Get(reflect.TypeOf(key))
	if err != nil {
		return nil, errtype.NewServiceWasNotFoundIntoContainerError(reflect.TypeOf(key))
	}
	service, ok := reflectService.Interface().(validatorinterface.Video)
	if !ok {
		return nil, errtype.NewTypesMismatchedServiceContainerError(reflect.TypeOf(reflectService), reflect.TypeOf(key))
	}
	return service, nil
}

func (s *ServiceContainer) GetVideoRepository() (repositoryinterface.Video, error) {
	key := (*repositoryinterface.Video)(nil)
	reflectService, err := s.Get(reflect.TypeOf(key))
	if err != nil {
		return nil, errtype.NewServiceWasNotFoundIntoContainerError(reflect.TypeOf(key))
	}
	service, ok := reflectService.Interface().(repositoryinterface.Video)
	if !ok {
		return nil, errtype.NewTypesMismatchedServiceContainerError(reflect.TypeOf(reflectService), reflect.TypeOf(key))
	}
	return service, nil
}

func (s *ServiceContainer) GetVideoCRUDService() (videoservice.CRUD, error) {
	key := (*videoservice.CRUD)(nil)
	reflectService, err := s.Get(reflect.TypeOf(key))
	if err != nil {
		return nil, errtype.NewServiceWasNotFoundIntoContainerError(reflect.TypeOf(key))
	}
	service, ok := reflectService.Interface().(videoservice.CRUD)
	if !ok {
		return nil, errtype.NewTypesMismatchedServiceContainerError(reflect.TypeOf(reflectService), reflect.TypeOf(key))
	}
	return service, nil
}

func (s *ServiceContainer) GetUserBuilder() (builderinterface.User, error) {
	key := (*builderinterface.User)(nil)
	reflectService, err := s.Get(reflect.TypeOf(key))
	if err != nil {
		return nil, errtype.NewServiceWasNotFoundIntoContainerError(reflect.TypeOf(key))
	}
	service, ok := reflectService.Interface().(builderinterface.User)
	if !ok {
		return nil, errtype.NewTypesMismatchedServiceContainerError(reflect.TypeOf(reflectService), reflect.TypeOf(key))
	}
	return service, nil
}

func (s *ServiceContainer) GetUserValidator() (validatorinterface.User, error) {
	key := (*validatorinterface.User)(nil)
	reflectService, err := s.Get(reflect.TypeOf(key))
	if err != nil {
		return nil, errtype.NewServiceWasNotFoundIntoContainerError(reflect.TypeOf(key))
	}
	service, ok := reflectService.Interface().(validatorinterface.User)
	if !ok {
		return nil, errtype.NewTypesMismatchedServiceContainerError(reflect.TypeOf(reflectService), reflect.TypeOf(key))
	}
	return service, nil
}

func (s *ServiceContainer) GetUserRepository() (repositoryinterface.User, error) {
	key := (*repositoryinterface.User)(nil)
	reflectService, err := s.Get(reflect.TypeOf(key))
	if err != nil {
		return nil, errtype.NewServiceWasNotFoundIntoContainerError(reflect.TypeOf(key))
	}
	service, ok := reflectService.Interface().(repositoryinterface.User)
	if !ok {
		return nil, errtype.NewTypesMismatchedServiceContainerError(reflect.TypeOf(reflectService), reflect.TypeOf(key))
	}
	return service, nil
}

func (s *ServiceContainer) GetUserCRUDService() (userservice.CRUD, error) {
	key := (*userservice.CRUD)(nil)
	reflectService, err := s.Get(reflect.TypeOf(key))
	if err != nil {
		return nil, errtype.NewServiceWasNotFoundIntoContainerError(reflect.TypeOf(key))
	}
	service, ok := reflectService.Interface().(userservice.CRUD)
	if !ok {
		return nil, errtype.NewTypesMismatchedServiceContainerError(reflect.TypeOf(reflectService), reflect.TypeOf(key))
	}
	return service, nil
}

func (s *ServiceContainer) GetAuthBuilder() (builderinterface.Auth, error) {
	key := (*builderinterface.Auth)(nil)
	reflectService, err := s.Get(reflect.TypeOf(key))
	if err != nil {
		return nil, errtype.NewServiceWasNotFoundIntoContainerError(reflect.TypeOf(key))
	}
	service, ok := reflectService.Interface().(builderinterface.Auth)
	if !ok {
		return nil, errtype.NewTypesMismatchedServiceContainerError(reflect.TypeOf(reflectService), reflect.TypeOf(key))
	}
	return service, nil
}

func (s *ServiceContainer) GetAuthValidator() (validatorinterface.Auth, error) {
	key := (*validatorinterface.Auth)(nil)
	reflectService, err := s.Get(reflect.TypeOf(key))
	if err != nil {
		return nil, errtype.NewServiceWasNotFoundIntoContainerError(reflect.TypeOf(key))
	}
	service, ok := reflectService.Interface().(validatorinterface.Auth)
	if !ok {
		return nil, errtype.NewTypesMismatchedServiceContainerError(reflect.TypeOf(reflectService), reflect.TypeOf(key))
	}
	return service, nil
}

func (s *ServiceContainer) GetAuthService() (authenticatorinterface.Authenticator, error) {
	key := (*authenticatorinterface.Authenticator)(nil)
	reflectService, err := s.Get(reflect.TypeOf(key))
	if err != nil {
		return nil, errtype.NewServiceWasNotFoundIntoContainerError(reflect.TypeOf(key))
	}
	service, ok := reflectService.Interface().(authenticatorinterface.Authenticator)
	if !ok {
		return nil, errtype.NewTypesMismatchedServiceContainerError(reflect.TypeOf(reflectService), reflect.TypeOf(key))
	}
	return service, nil
}

func (s *ServiceContainer) GetLoggerService() (loggerinterface.Logger, error) {
	key := (*loggerinterface.Logger)(nil)
	reflectService, err := s.Get(reflect.TypeOf(key))
	if err != nil {
		return nil, errtype.NewServiceWasNotFoundIntoContainerError(reflect.TypeOf(key))
	}
	service, ok := reflectService.Interface().(loggerinterface.Logger)
	if !ok {
		return nil, errtype.NewTypesMismatchedServiceContainerError(reflect.TypeOf(reflectService), reflect.TypeOf(key))
	}
	return service, nil
}

func (s *ServiceContainer) GetCacheService() (cacherinterface.Cacher, error) {
	key := (*cacherinterface.Cacher)(nil)
	reflectService, err := s.Get(reflect.TypeOf(key))
	if err != nil {
		return nil, errtype.NewServiceWasNotFoundIntoContainerError(reflect.TypeOf(key))
	}
	service, ok := reflectService.Interface().(cacherinterface.Cacher)
	if !ok {
		return nil, errtype.NewTypesMismatchedServiceContainerError(reflect.TypeOf(reflectService), reflect.TypeOf(key))
	}
	return service, nil
}

func (s *ServiceContainer) GetFileStorageService() (fileinterface.Storage, error) {
	key := (*fileinterface.Storage)(nil)
	reflectService, err := s.Get(reflect.TypeOf(key))
	if err != nil {
		return nil, errtype.NewServiceWasNotFoundIntoContainerError(reflect.TypeOf(key))
	}
	service, ok := reflectService.Interface().(fileinterface.Storage)
	if !ok {
		return nil, errtype.NewTypesMismatchedServiceContainerError(reflect.TypeOf(reflectService), reflect.TypeOf(key))
	}
	return service, nil
}

func (s *ServiceContainer) GetFileNameComputerService() (fileinterface.NameComputer, error) {
	key := (*fileinterface.NameComputer)(nil)
	reflectService, err := s.Get(reflect.TypeOf(key))
	if err != nil {
		return nil, errtype.NewServiceWasNotFoundIntoContainerError(reflect.TypeOf(key))
	}
	service, ok := reflectService.Interface().(fileinterface.NameComputer)
	if !ok {
		return nil, errtype.NewTypesMismatchedServiceContainerError(reflect.TypeOf(reflectService), reflect.TypeOf(key))
	}
	return service, nil
}

func (s *ServiceContainer) GetFileUploaderService() (uploaderinterface.Uploader, error) {
	key := (*uploaderinterface.Uploader)(nil)
	reflectService, err := s.Get(reflect.TypeOf(key))
	if err != nil {
		return nil, errtype.NewServiceWasNotFoundIntoContainerError(reflect.TypeOf(key))
	}
	service, ok := reflectService.Interface().(uploaderinterface.Uploader)
	if !ok {
		return nil, errtype.NewTypesMismatchedServiceContainerError(reflect.TypeOf(reflectService), reflect.TypeOf(key))
	}
	return service, nil
}

func (s *ServiceContainer) GetFileReaderService() (readerinterface.FileReader, error) {
	key := (*readerinterface.FileReader)(nil)
	reflectService, err := s.Get(reflect.TypeOf(key))
	if err != nil {
		return nil, errtype.NewServiceWasNotFoundIntoContainerError(reflect.TypeOf(key))
	}
	service, ok := reflectService.Interface().(readerinterface.FileReader)
	if !ok {
		return nil, errtype.NewTypesMismatchedServiceContainerError(reflect.TypeOf(reflectService), reflect.TypeOf(key))
	}
	return service, nil
}

func (s *ServiceContainer) GetRequestParametersExtractorService() (extractorinterface.RequestParams, error) {
	key := (*extractorinterface.RequestParams)(nil)
	reflectService, err := s.Get(reflect.TypeOf(key))
	if err != nil {
		return nil, errtype.NewServiceWasNotFoundIntoContainerError(reflect.TypeOf(key))
	}
	service, ok := reflectService.Interface().(extractorinterface.RequestParams)
	if !ok {
		return nil, errtype.NewTypesMismatchedServiceContainerError(reflect.TypeOf(reflectService), reflect.TypeOf(key))
	}
	return service, nil
}

func (s *ServiceContainer) GetResponderService() (responseinterface.Responder, error) {
	key := (*responseinterface.Responder)(nil)
	reflectService, err := s.Get(reflect.TypeOf(key))
	if err != nil {
		return nil, errtype.NewServiceWasNotFoundIntoContainerError(reflect.TypeOf(key))
	}
	service, ok := reflectService.Interface().(responseinterface.Responder)
	if !ok {
		return nil, errtype.NewTypesMismatchedServiceContainerError(reflect.TypeOf(reflectService), reflect.TypeOf(key))
	}
	return service, nil
}

func (s *ServiceContainer) GetPasswordHasherService() (securityinterface.PasswordHasher, error) {
	key := (*securityinterface.PasswordHasher)(nil)
	reflectService, err := s.Get(reflect.TypeOf(key))
	if err != nil {
		return nil, errtype.NewServiceWasNotFoundIntoContainerError(reflect.TypeOf(key))
	}
	service, ok := reflectService.Interface().(securityinterface.PasswordHasher)
	if !ok {
		return nil, errtype.NewTypesMismatchedServiceContainerError(reflect.TypeOf(reflectService), reflect.TypeOf(key))
	}
	return service, nil
}

func (s *ServiceContainer) GetTokenizerService() (tokenizerinterface.Tokenizer, error) {
	key := (*tokenizerinterface.Tokenizer)(nil)
	reflectService, err := s.Get(reflect.TypeOf(key))
	if err != nil {
		return nil, errtype.NewServiceWasNotFoundIntoContainerError(reflect.TypeOf(key))
	}
	service, ok := reflectService.Interface().(tokenizerinterface.Tokenizer)
	if !ok {
		return nil, errtype.NewTypesMismatchedServiceContainerError(reflect.TypeOf(reflectService), reflect.TypeOf(key))
	}
	return service, nil
}

func (s *ServiceContainer) GetWebSocketCommunicatorService() (protointerface.Communicator, error) {
	key := (*protointerface.Communicator)(nil)
	reflectService, err := s.Get(reflect.TypeOf(key))
	if err != nil {
		return nil, errtype.NewServiceWasNotFoundIntoContainerError(reflect.TypeOf(key))
	}
	service, ok := reflectService.Interface().(protointerface.Communicator)
	if !ok {
		return nil, errtype.NewTypesMismatchedServiceContainerError(reflect.TypeOf(reflectService), reflect.TypeOf(key))
	}
	return service, nil
}

func (s *ServiceContainer) GetWebSocketHandlerStrategies() ([]strategyinterface.ActionStrategy, error) {
	key := (*[]strategyinterface.ActionStrategy)(nil)
	reflectService, err := s.Get(reflect.TypeOf(key))
	if err != nil {
		return nil, errtype.NewServiceWasNotFoundIntoContainerError(reflect.TypeOf(key))
	}
	service, ok := reflectService.Interface().([]strategyinterface.ActionStrategy)
	if !ok {
		return nil, errtype.NewTypesMismatchedServiceContainerError(reflect.TypeOf(reflectService), reflect.TypeOf(key))
	}
	return service, nil
}

func (s *ServiceContainer) GetWebSocketListener() (listenerinterface.ActionsListener, error) {
	key := (*listenerinterface.ActionsListener)(nil)
	reflectService, err := s.Get(reflect.TypeOf(key))
	if err != nil {
		return nil, errtype.NewServiceWasNotFoundIntoContainerError(reflect.TypeOf(key))
	}
	service, ok := reflectService.Interface().(listenerinterface.ActionsListener)
	if !ok {
		return nil, errtype.NewTypesMismatchedServiceContainerError(reflect.TypeOf(reflectService), reflect.TypeOf(key))
	}
	return service, nil
}

func (s *ServiceContainer) GetWebSocketHandler() (handlerinterface.ActionsHandler, error) {
	key := (*handlerinterface.ActionsHandler)(nil)
	reflectService, err := s.Get(reflect.TypeOf(key))
	if err != nil {
		return nil, errtype.NewServiceWasNotFoundIntoContainerError(reflect.TypeOf(key))
	}
	service, ok := reflectService.Interface().(handlerinterface.ActionsHandler)
	if !ok {
		return nil, errtype.NewTypesMismatchedServiceContainerError(reflect.TypeOf(reflectService), reflect.TypeOf(key))
	}
	return service, nil
}

func (s *ServiceContainer) GetCodecsDetectorService() (detectorinterface.Codecs, error) {
	key := (*detectorinterface.Codecs)(nil)
	reflectService, err := s.Get(reflect.TypeOf(key))
	if err != nil {
		return nil, errtype.NewServiceWasNotFoundIntoContainerError(reflect.TypeOf(key))
	}
	service, ok := reflectService.Interface().(detectorinterface.Codecs)
	if !ok {
		return nil, errtype.NewTypesMismatchedServiceContainerError(reflect.TypeOf(reflectService), reflect.TypeOf(key))
	}
	return service, nil
}

func (s *ServiceContainer) GetStreamingService() (streamerinterface.Streamer, error) {
	key := (*streamerinterface.Streamer)(nil)
	reflectService, err := s.Get(reflect.TypeOf(key))
	if err != nil {
		return nil, errtype.NewServiceWasNotFoundIntoContainerError(reflect.TypeOf(key))
	}
	service, ok := reflectService.Interface().(streamerinterface.Streamer)
	if !ok {
		return nil, errtype.NewTypesMismatchedServiceContainerError(reflect.TypeOf(reflectService), reflect.TypeOf(key))
	}
	return service, nil
}
