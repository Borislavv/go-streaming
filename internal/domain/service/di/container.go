package di

import (
	"context"
	"github.com/Borislavv/video-streaming/internal/app"
	builderinterface "github.com/Borislavv/video-streaming/internal/domain/builder/interface"
	"github.com/Borislavv/video-streaming/internal/domain/errors"
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
	"github.com/Borislavv/video-streaming/internal/infrastructure/di"
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

type ServiceContainerManager struct {
	diinterface.Container
}

func NewServiceContainerManager() *ServiceContainerManager {
	return &ServiceContainerManager{
		Container: di.NewServiceContainer(),
	}
}

func (s *ServiceContainerManager) GetConfig() (*app.Config, error) {
	key := (*app.Config)(nil)
	service, err := s.Get(reflect.TypeOf(key))
	if err != nil {
		return nil, errors.NewServiceWasNotFoundIntoContainerError(reflect.TypeOf(key))
	}
	cfg, ok := service.Interface().(*app.Config)
	if !ok {
		return nil, errors.NewTypesMismatchedServiceContainerError(reflect.TypeOf(service), reflect.TypeOf(key))
	}
	return cfg, nil
}

func (s *ServiceContainerManager) GetCtx() (context.Context, error) {
	key := (*context.Context)(nil)
	service, err := s.Get(reflect.TypeOf(key))
	if err != nil {
		return nil, errors.NewServiceWasNotFoundIntoContainerError(reflect.TypeOf(key))
	}
	ctx, ok := service.Interface().(context.Context)
	if !ok {
		return nil, errors.NewTypesMismatchedServiceContainerError(reflect.TypeOf(service), reflect.TypeOf(key))
	}
	return ctx, nil
}

func (s *ServiceContainerManager) GetCancelFunc() (context.CancelFunc, error) {
	key := (*context.CancelFunc)(nil)
	service, err := s.Get(reflect.TypeOf(key))
	if err != nil {
		return nil, errors.NewServiceWasNotFoundIntoContainerError(reflect.TypeOf(key))
	}
	cancel, ok := service.Interface().(context.CancelFunc)
	if !ok {
		return nil, errors.NewTypesMismatchedServiceContainerError(reflect.TypeOf(service), reflect.TypeOf(key))
	}
	return cancel, nil
}

func (s *ServiceContainerManager) GetMongoDatabase() (*mongo.Database, error) {
	key := (*mongo.Database)(nil)
	service, err := s.Get(reflect.TypeOf(key))
	if err != nil {
		return nil, errors.NewServiceWasNotFoundIntoContainerError(reflect.TypeOf(key))
	}
	database, ok := service.Interface().(*mongo.Database)
	if !ok {
		return nil, errors.NewTypesMismatchedServiceContainerError(reflect.TypeOf(service), reflect.TypeOf(key))
	}
	return database, nil
}

func (s *ServiceContainerManager) GetResourceMongoRepository() (mongodbinterface.Resource, error) {
	key := (*mongodbinterface.Resource)(nil)
	service, err := s.Get(reflect.TypeOf(key))
	if err != nil {
		return nil, errors.NewServiceWasNotFoundIntoContainerError(reflect.TypeOf(key))
	}
	repo, ok := service.Interface().(mongodbinterface.Resource)
	if !ok {
		return nil, errors.NewTypesMismatchedServiceContainerError(reflect.TypeOf(service), reflect.TypeOf(key))
	}
	return repo, nil
}

func (s *ServiceContainerManager) GetVideoMongoRepository() (mongodbinterface.Video, error) {
	key := (*mongodbinterface.Video)(nil)
	service, err := s.Get(reflect.TypeOf(key))
	if err != nil {
		return nil, errors.NewServiceWasNotFoundIntoContainerError(reflect.TypeOf(key))
	}
	repo, ok := service.Interface().(mongodbinterface.Video)
	if !ok {
		return nil, errors.NewTypesMismatchedServiceContainerError(reflect.TypeOf(service), reflect.TypeOf(key))
	}
	return repo, nil
}

func (s *ServiceContainerManager) GetUserMongoRepository() (mongodbinterface.User, error) {
	key := (*mongodbinterface.User)(nil)
	service, err := s.Get(reflect.TypeOf(key))
	if err != nil {
		return nil, errors.NewServiceWasNotFoundIntoContainerError(reflect.TypeOf(key))
	}
	repo, ok := service.Interface().(mongodbinterface.User)
	if !ok {
		return nil, errors.NewTypesMismatchedServiceContainerError(reflect.TypeOf(service), reflect.TypeOf(key))
	}
	return repo, nil
}

func (s *ServiceContainerManager) GetBlockedTokenMongoRepository() (mongodbinterface.BlockedToken, error) {
	key := (*mongodbinterface.BlockedToken)(nil)
	service, err := s.Get(reflect.TypeOf(key))
	if err != nil {
		return nil, errors.NewServiceWasNotFoundIntoContainerError(reflect.TypeOf(key))
	}
	repo, ok := service.Interface().(mongodbinterface.BlockedToken)
	if !ok {
		return nil, errors.NewTypesMismatchedServiceContainerError(reflect.TypeOf(service), reflect.TypeOf(key))
	}
	return repo, nil
}

func (s *ServiceContainerManager) GetResourceCacheRepository() (cacheinterface.Resource, error) {
	key := (*cacheinterface.Resource)(nil)
	service, err := s.Get(reflect.TypeOf(key))
	if err != nil {
		return nil, errors.NewServiceWasNotFoundIntoContainerError(reflect.TypeOf(key))
	}
	repo, ok := service.Interface().(cacheinterface.Resource)
	if !ok {
		return nil, errors.NewTypesMismatchedServiceContainerError(reflect.TypeOf(service), reflect.TypeOf(key))
	}
	return repo, nil
}

func (s *ServiceContainerManager) GetVideoCacheRepository() (cacheinterface.Video, error) {
	key := (*cacheinterface.Video)(nil)
	service, err := s.Get(reflect.TypeOf(key))
	if err != nil {
		return nil, errors.NewServiceWasNotFoundIntoContainerError(reflect.TypeOf(key))
	}
	repo, ok := service.Interface().(cacheinterface.Video)
	if !ok {
		return nil, errors.NewTypesMismatchedServiceContainerError(reflect.TypeOf(service), reflect.TypeOf(key))
	}
	return repo, nil
}

func (s *ServiceContainerManager) GetUserCacheRepository() (cacheinterface.User, error) {
	key := (*cacheinterface.User)(nil)
	service, err := s.Get(reflect.TypeOf(key))
	if err != nil {
		return nil, errors.NewServiceWasNotFoundIntoContainerError(reflect.TypeOf(key))
	}
	repo, ok := service.Interface().(cacheinterface.User)
	if !ok {
		return nil, errors.NewTypesMismatchedServiceContainerError(reflect.TypeOf(service), reflect.TypeOf(key))
	}
	return repo, nil
}

func (s *ServiceContainerManager) GetAccessService() (accessorinterface.Accessor, error) {
	key := (*accessorinterface.Accessor)(nil)
	service, err := s.Get(reflect.TypeOf(key))
	if err != nil {
		return nil, errors.NewServiceWasNotFoundIntoContainerError(reflect.TypeOf(key))
	}
	accessService, ok := service.Interface().(accessorinterface.Accessor)
	if !ok {
		return nil, errors.NewTypesMismatchedServiceContainerError(reflect.TypeOf(service), reflect.TypeOf(key))
	}
	return accessService, nil
}

func (s *ServiceContainerManager) GetResourceBuilder() (builderinterface.Resource, error) {
	key := (*builderinterface.Resource)(nil)
	service, err := s.Get(reflect.TypeOf(key))
	if err != nil {
		return nil, errors.NewServiceWasNotFoundIntoContainerError(reflect.TypeOf(key))
	}
	builderService, ok := service.Interface().(builderinterface.Resource)
	if !ok {
		return nil, errors.NewTypesMismatchedServiceContainerError(reflect.TypeOf(service), reflect.TypeOf(key))
	}
	return builderService, nil
}

func (s *ServiceContainerManager) GetResourceValidator() (validatorinterface.Resource, error) {
	key := (*validatorinterface.Resource)(nil)
	service, err := s.Get(reflect.TypeOf(key))
	if err != nil {
		return nil, errors.NewServiceWasNotFoundIntoContainerError(reflect.TypeOf(key))
	}
	validatorService, ok := service.Interface().(validatorinterface.Resource)
	if !ok {
		return nil, errors.NewTypesMismatchedServiceContainerError(reflect.TypeOf(service), reflect.TypeOf(key))
	}
	return validatorService, nil
}

func (s *ServiceContainerManager) GetResourceRepository() (repositoryinterface.Resource, error) {
	key := (*repositoryinterface.Resource)(nil)
	service, err := s.Get(reflect.TypeOf(key))
	if err != nil {
		return nil, errors.NewServiceWasNotFoundIntoContainerError(reflect.TypeOf(key))
	}
	repo, ok := service.Interface().(repositoryinterface.Resource)
	if !ok {
		return nil, errors.NewTypesMismatchedServiceContainerError(reflect.TypeOf(service), reflect.TypeOf(key))
	}
	return repo, nil
}

func (s *ServiceContainerManager) GetResourceCRUDService() (resourceservice.CRUD, error) {
	key := (*resourceservice.CRUD)(nil)
	reflectService, err := s.Get(reflect.TypeOf(key))
	if err != nil {
		return nil, errors.NewServiceWasNotFoundIntoContainerError(reflect.TypeOf(key))
	}
	service, ok := reflectService.Interface().(resourceservice.CRUD)
	if !ok {
		return nil, errors.NewTypesMismatchedServiceContainerError(reflect.TypeOf(reflectService), reflect.TypeOf(key))
	}
	return service, nil
}

func (s *ServiceContainerManager) GetBlockedTokenRepository() (repositoryinterface.BlockedToken, error) {
	key := (*repositoryinterface.BlockedToken)(nil)
	reflectService, err := s.Get(reflect.TypeOf(key))
	if err != nil {
		return nil, errors.NewServiceWasNotFoundIntoContainerError(reflect.TypeOf(key))
	}
	service, ok := reflectService.Interface().(repositoryinterface.BlockedToken)
	if !ok {
		return nil, errors.NewTypesMismatchedServiceContainerError(reflect.TypeOf(reflectService), reflect.TypeOf(key))
	}
	return service, nil
}

func (s *ServiceContainerManager) GetVideoBuilder() (builderinterface.Video, error) {
	key := (*builderinterface.Video)(nil)
	reflectService, err := s.Get(reflect.TypeOf(key))
	if err != nil {
		return nil, errors.NewServiceWasNotFoundIntoContainerError(reflect.TypeOf(key))
	}
	service, ok := reflectService.Interface().(builderinterface.Video)
	if !ok {
		return nil, errors.NewTypesMismatchedServiceContainerError(reflect.TypeOf(reflectService), reflect.TypeOf(key))
	}
	return service, nil
}

func (s *ServiceContainerManager) GetVideoValidator() (validatorinterface.Video, error) {
	key := (*validatorinterface.Video)(nil)
	reflectService, err := s.Get(reflect.TypeOf(key))
	if err != nil {
		return nil, errors.NewServiceWasNotFoundIntoContainerError(reflect.TypeOf(key))
	}
	service, ok := reflectService.Interface().(validatorinterface.Video)
	if !ok {
		return nil, errors.NewTypesMismatchedServiceContainerError(reflect.TypeOf(reflectService), reflect.TypeOf(key))
	}
	return service, nil
}

func (s *ServiceContainerManager) GetVideoRepository() (repositoryinterface.Video, error) {
	key := (*repositoryinterface.Video)(nil)
	reflectService, err := s.Get(reflect.TypeOf(key))
	if err != nil {
		return nil, errors.NewServiceWasNotFoundIntoContainerError(reflect.TypeOf(key))
	}
	service, ok := reflectService.Interface().(repositoryinterface.Video)
	if !ok {
		return nil, errors.NewTypesMismatchedServiceContainerError(reflect.TypeOf(reflectService), reflect.TypeOf(key))
	}
	return service, nil
}

func (s *ServiceContainerManager) GetVideoCRUDService() (videoservice.CRUD, error) {
	key := (*videoservice.CRUD)(nil)
	reflectService, err := s.Get(reflect.TypeOf(key))
	if err != nil {
		return nil, errors.NewServiceWasNotFoundIntoContainerError(reflect.TypeOf(key))
	}
	service, ok := reflectService.Interface().(videoservice.CRUD)
	if !ok {
		return nil, errors.NewTypesMismatchedServiceContainerError(reflect.TypeOf(reflectService), reflect.TypeOf(key))
	}
	return service, nil
}

func (s *ServiceContainerManager) GetUserBuilder() (builderinterface.User, error) {
	key := (*builderinterface.User)(nil)
	reflectService, err := s.Get(reflect.TypeOf(key))
	if err != nil {
		return nil, errors.NewServiceWasNotFoundIntoContainerError(reflect.TypeOf(key))
	}
	service, ok := reflectService.Interface().(builderinterface.User)
	if !ok {
		return nil, errors.NewTypesMismatchedServiceContainerError(reflect.TypeOf(reflectService), reflect.TypeOf(key))
	}
	return service, nil
}

func (s *ServiceContainerManager) GetUserValidator() (validatorinterface.User, error) {
	key := (*validatorinterface.User)(nil)
	reflectService, err := s.Get(reflect.TypeOf(key))
	if err != nil {
		return nil, errors.NewServiceWasNotFoundIntoContainerError(reflect.TypeOf(key))
	}
	service, ok := reflectService.Interface().(validatorinterface.User)
	if !ok {
		return nil, errors.NewTypesMismatchedServiceContainerError(reflect.TypeOf(reflectService), reflect.TypeOf(key))
	}
	return service, nil
}

func (s *ServiceContainerManager) GetUserRepository() (repositoryinterface.User, error) {
	key := (*repositoryinterface.User)(nil)
	reflectService, err := s.Get(reflect.TypeOf(key))
	if err != nil {
		return nil, errors.NewServiceWasNotFoundIntoContainerError(reflect.TypeOf(key))
	}
	service, ok := reflectService.Interface().(repositoryinterface.User)
	if !ok {
		return nil, errors.NewTypesMismatchedServiceContainerError(reflect.TypeOf(reflectService), reflect.TypeOf(key))
	}
	return service, nil
}

func (s *ServiceContainerManager) GetUserCRUDService() (userservice.CRUD, error) {
	key := (*userservice.CRUD)(nil)
	reflectService, err := s.Get(reflect.TypeOf(key))
	if err != nil {
		return nil, errors.NewServiceWasNotFoundIntoContainerError(reflect.TypeOf(key))
	}
	service, ok := reflectService.Interface().(userservice.CRUD)
	if !ok {
		return nil, errors.NewTypesMismatchedServiceContainerError(reflect.TypeOf(reflectService), reflect.TypeOf(key))
	}
	return service, nil
}

func (s *ServiceContainerManager) GetAuthBuilder() (builderinterface.Auth, error) {
	key := (*builderinterface.Auth)(nil)
	reflectService, err := s.Get(reflect.TypeOf(key))
	if err != nil {
		return nil, errors.NewServiceWasNotFoundIntoContainerError(reflect.TypeOf(key))
	}
	service, ok := reflectService.Interface().(builderinterface.Auth)
	if !ok {
		return nil, errors.NewTypesMismatchedServiceContainerError(reflect.TypeOf(reflectService), reflect.TypeOf(key))
	}
	return service, nil
}

func (s *ServiceContainerManager) GetAuthValidator() (validatorinterface.Auth, error) {
	key := (*validatorinterface.Auth)(nil)
	reflectService, err := s.Get(reflect.TypeOf(key))
	if err != nil {
		return nil, errors.NewServiceWasNotFoundIntoContainerError(reflect.TypeOf(key))
	}
	service, ok := reflectService.Interface().(validatorinterface.Auth)
	if !ok {
		return nil, errors.NewTypesMismatchedServiceContainerError(reflect.TypeOf(reflectService), reflect.TypeOf(key))
	}
	return service, nil
}

func (s *ServiceContainerManager) GetAuthService() (authenticatorinterface.Authenticator, error) {
	key := (*authenticatorinterface.Authenticator)(nil)
	reflectService, err := s.Get(reflect.TypeOf(key))
	if err != nil {
		return nil, errors.NewServiceWasNotFoundIntoContainerError(reflect.TypeOf(key))
	}
	service, ok := reflectService.Interface().(authenticatorinterface.Authenticator)
	if !ok {
		return nil, errors.NewTypesMismatchedServiceContainerError(reflect.TypeOf(reflectService), reflect.TypeOf(key))
	}
	return service, nil
}

func (s *ServiceContainerManager) GetLoggerService() (loggerinterface.Logger, error) {
	key := (*loggerinterface.Logger)(nil)
	reflectService, err := s.Get(reflect.TypeOf(key))
	if err != nil {
		return nil, errors.NewServiceWasNotFoundIntoContainerError(reflect.TypeOf(key))
	}
	service, ok := reflectService.Interface().(loggerinterface.Logger)
	if !ok {
		return nil, errors.NewTypesMismatchedServiceContainerError(reflect.TypeOf(reflectService), reflect.TypeOf(key))
	}
	return service, nil
}

func (s *ServiceContainerManager) GetCacheService() (cacherinterface.Cacher, error) {
	key := (*cacherinterface.Cacher)(nil)
	reflectService, err := s.Get(reflect.TypeOf(key))
	if err != nil {
		return nil, errors.NewServiceWasNotFoundIntoContainerError(reflect.TypeOf(key))
	}
	service, ok := reflectService.Interface().(cacherinterface.Cacher)
	if !ok {
		return nil, errors.NewTypesMismatchedServiceContainerError(reflect.TypeOf(reflectService), reflect.TypeOf(key))
	}
	return service, nil
}

func (s *ServiceContainerManager) GetFileStorageService() (fileinterface.Storage, error) {
	key := (*fileinterface.Storage)(nil)
	reflectService, err := s.Get(reflect.TypeOf(key))
	if err != nil {
		return nil, errors.NewServiceWasNotFoundIntoContainerError(reflect.TypeOf(key))
	}
	service, ok := reflectService.Interface().(fileinterface.Storage)
	if !ok {
		return nil, errors.NewTypesMismatchedServiceContainerError(reflect.TypeOf(reflectService), reflect.TypeOf(key))
	}
	return service, nil
}

func (s *ServiceContainerManager) GetFileNameComputerService() (fileinterface.NameComputer, error) {
	key := (*fileinterface.NameComputer)(nil)
	reflectService, err := s.Get(reflect.TypeOf(key))
	if err != nil {
		return nil, errors.NewServiceWasNotFoundIntoContainerError(reflect.TypeOf(key))
	}
	service, ok := reflectService.Interface().(fileinterface.NameComputer)
	if !ok {
		return nil, errors.NewTypesMismatchedServiceContainerError(reflect.TypeOf(reflectService), reflect.TypeOf(key))
	}
	return service, nil
}

func (s *ServiceContainerManager) GetFileUploaderService() (uploaderinterface.Uploader, error) {
	key := (*uploaderinterface.Uploader)(nil)
	reflectService, err := s.Get(reflect.TypeOf(key))
	if err != nil {
		return nil, errors.NewServiceWasNotFoundIntoContainerError(reflect.TypeOf(key))
	}
	service, ok := reflectService.Interface().(uploaderinterface.Uploader)
	if !ok {
		return nil, errors.NewTypesMismatchedServiceContainerError(reflect.TypeOf(reflectService), reflect.TypeOf(key))
	}
	return service, nil
}

func (s *ServiceContainerManager) GetFileReaderService() (readerinterface.FileReader, error) {
	key := (*readerinterface.FileReader)(nil)
	reflectService, err := s.Get(reflect.TypeOf(key))
	if err != nil {
		return nil, errors.NewServiceWasNotFoundIntoContainerError(reflect.TypeOf(key))
	}
	service, ok := reflectService.Interface().(readerinterface.FileReader)
	if !ok {
		return nil, errors.NewTypesMismatchedServiceContainerError(reflect.TypeOf(reflectService), reflect.TypeOf(key))
	}
	return service, nil
}

func (s *ServiceContainerManager) GetRequestParametersExtractorService() (extractorinterface.RequestParams, error) {
	key := (*extractorinterface.RequestParams)(nil)
	reflectService, err := s.Get(reflect.TypeOf(key))
	if err != nil {
		return nil, errors.NewServiceWasNotFoundIntoContainerError(reflect.TypeOf(key))
	}
	service, ok := reflectService.Interface().(extractorinterface.RequestParams)
	if !ok {
		return nil, errors.NewTypesMismatchedServiceContainerError(reflect.TypeOf(reflectService), reflect.TypeOf(key))
	}
	return service, nil
}

func (s *ServiceContainerManager) GetResponderService() (responseinterface.Responder, error) {
	key := (*responseinterface.Responder)(nil)
	reflectService, err := s.Get(reflect.TypeOf(key))
	if err != nil {
		return nil, errors.NewServiceWasNotFoundIntoContainerError(reflect.TypeOf(key))
	}
	service, ok := reflectService.Interface().(responseinterface.Responder)
	if !ok {
		return nil, errors.NewTypesMismatchedServiceContainerError(reflect.TypeOf(reflectService), reflect.TypeOf(key))
	}
	return service, nil
}

func (s *ServiceContainerManager) GetPasswordHasherService() (securityinterface.PasswordHasher, error) {
	key := (*securityinterface.PasswordHasher)(nil)
	reflectService, err := s.Get(reflect.TypeOf(key))
	if err != nil {
		return nil, errors.NewServiceWasNotFoundIntoContainerError(reflect.TypeOf(key))
	}
	service, ok := reflectService.Interface().(securityinterface.PasswordHasher)
	if !ok {
		return nil, errors.NewTypesMismatchedServiceContainerError(reflect.TypeOf(reflectService), reflect.TypeOf(key))
	}
	return service, nil
}

func (s *ServiceContainerManager) GetTokenizerService() (tokenizerinterface.Tokenizer, error) {
	key := (*tokenizerinterface.Tokenizer)(nil)
	reflectService, err := s.Get(reflect.TypeOf(key))
	if err != nil {
		return nil, errors.NewServiceWasNotFoundIntoContainerError(reflect.TypeOf(key))
	}
	service, ok := reflectService.Interface().(tokenizerinterface.Tokenizer)
	if !ok {
		return nil, errors.NewTypesMismatchedServiceContainerError(reflect.TypeOf(reflectService), reflect.TypeOf(key))
	}
	return service, nil
}

func (s *ServiceContainerManager) GetWebSocketCommunicatorService() (protointerface.Communicator, error) {
	key := (*protointerface.Communicator)(nil)
	reflectService, err := s.Get(reflect.TypeOf(key))
	if err != nil {
		return nil, errors.NewServiceWasNotFoundIntoContainerError(reflect.TypeOf(key))
	}
	service, ok := reflectService.Interface().(protointerface.Communicator)
	if !ok {
		return nil, errors.NewTypesMismatchedServiceContainerError(reflect.TypeOf(reflectService), reflect.TypeOf(key))
	}
	return service, nil
}

func (s *ServiceContainerManager) GetWebSocketHandlerStrategies() ([]strategyinterface.ActionStrategy, error) {
	key := (*[]strategyinterface.ActionStrategy)(nil)
	reflectService, err := s.Get(reflect.TypeOf(key))
	if err != nil {
		return nil, errors.NewServiceWasNotFoundIntoContainerError(reflect.TypeOf(key))
	}
	service, ok := reflectService.Interface().([]strategyinterface.ActionStrategy)
	if !ok {
		return nil, errors.NewTypesMismatchedServiceContainerError(reflect.TypeOf(reflectService), reflect.TypeOf(key))
	}
	return service, nil
}

func (s *ServiceContainerManager) GetWebSocketListener() (listenerinterface.ActionsListener, error) {
	key := (*listenerinterface.ActionsListener)(nil)
	reflectService, err := s.Get(reflect.TypeOf(key))
	if err != nil {
		return nil, errors.NewServiceWasNotFoundIntoContainerError(reflect.TypeOf(key))
	}
	service, ok := reflectService.Interface().(listenerinterface.ActionsListener)
	if !ok {
		return nil, errors.NewTypesMismatchedServiceContainerError(reflect.TypeOf(reflectService), reflect.TypeOf(key))
	}
	return service, nil
}

func (s *ServiceContainerManager) GetWebSocketHandler() (handlerinterface.ActionsHandler, error) {
	key := (*handlerinterface.ActionsHandler)(nil)
	reflectService, err := s.Get(reflect.TypeOf(key))
	if err != nil {
		return nil, errors.NewServiceWasNotFoundIntoContainerError(reflect.TypeOf(key))
	}
	service, ok := reflectService.Interface().(handlerinterface.ActionsHandler)
	if !ok {
		return nil, errors.NewTypesMismatchedServiceContainerError(reflect.TypeOf(reflectService), reflect.TypeOf(key))
	}
	return service, nil
}

func (s *ServiceContainerManager) GetCodecsDetectorService() (detectorinterface.Codecs, error) {
	key := (*detectorinterface.Codecs)(nil)
	reflectService, err := s.Get(reflect.TypeOf(key))
	if err != nil {
		return nil, errors.NewServiceWasNotFoundIntoContainerError(reflect.TypeOf(key))
	}
	service, ok := reflectService.Interface().(detectorinterface.Codecs)
	if !ok {
		return nil, errors.NewTypesMismatchedServiceContainerError(reflect.TypeOf(reflectService), reflect.TypeOf(key))
	}
	return service, nil
}

func (s *ServiceContainerManager) GetStreamingService() (streamerinterface.Streamer, error) {
	key := (*streamerinterface.Streamer)(nil)
	reflectService, err := s.Get(reflect.TypeOf(key))
	if err != nil {
		return nil, errors.NewServiceWasNotFoundIntoContainerError(reflect.TypeOf(key))
	}
	service, ok := reflectService.Interface().(streamerinterface.Streamer)
	if !ok {
		return nil, errors.NewTypesMismatchedServiceContainerError(reflect.TypeOf(reflectService), reflect.TypeOf(key))
	}
	return service, nil
}
