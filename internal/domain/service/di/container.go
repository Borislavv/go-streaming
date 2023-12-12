package di

import (
	"context"
	"github.com/Borislavv/video-streaming/internal/app"
	builder_interface "github.com/Borislavv/video-streaming/internal/domain/builder/interface"
	"github.com/Borislavv/video-streaming/internal/domain/errors"
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
	"github.com/Borislavv/video-streaming/internal/infrastructure/di"
	di_interface "github.com/Borislavv/video-streaming/internal/infrastructure/di/interface"
	cache_interface "github.com/Borislavv/video-streaming/internal/infrastructure/repository/storage/cache/interface"
	mongodb_interface "github.com/Borislavv/video-streaming/internal/infrastructure/repository/storage/mongodb/interface"
	detector_interface "github.com/Borislavv/video-streaming/internal/infrastructure/service/detector/interface"
	reader_interface "github.com/Borislavv/video-streaming/internal/infrastructure/service/reader/interface"
	handler_interface "github.com/Borislavv/video-streaming/internal/infrastructure/service/streamer/action/handler/interface"
	strategy_interface "github.com/Borislavv/video-streaming/internal/infrastructure/service/streamer/action/handler/strategy/interface"
	listener_interface "github.com/Borislavv/video-streaming/internal/infrastructure/service/streamer/action/listener/interface"
	streamer_interface "github.com/Borislavv/video-streaming/internal/infrastructure/service/streamer/interface"
	proto_interface "github.com/Borislavv/video-streaming/internal/infrastructure/service/streamer/proto/interface"
	file_interface "github.com/Borislavv/video-streaming/internal/infrastructure/service/uploader/file/interface"
	"go.mongodb.org/mongo-driver/mongo"
	"reflect"
)

type ServiceContainerManager struct {
	di_interface.Container
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

func (s *ServiceContainerManager) GetResourceMongoRepository() (mongodb_interface.Resource, error) {
	key := (*mongodb_interface.Resource)(nil)
	service, err := s.Get(reflect.TypeOf(key))
	if err != nil {
		return nil, errors.NewServiceWasNotFoundIntoContainerError(reflect.TypeOf(key))
	}
	repo, ok := service.Interface().(mongodb_interface.Resource)
	if !ok {
		return nil, errors.NewTypesMismatchedServiceContainerError(reflect.TypeOf(service), reflect.TypeOf(key))
	}
	return repo, nil
}

func (s *ServiceContainerManager) GetVideoMongoRepository() (mongodb_interface.Video, error) {
	key := (*mongodb_interface.Video)(nil)
	service, err := s.Get(reflect.TypeOf(key))
	if err != nil {
		return nil, errors.NewServiceWasNotFoundIntoContainerError(reflect.TypeOf(key))
	}
	repo, ok := service.Interface().(mongodb_interface.Video)
	if !ok {
		return nil, errors.NewTypesMismatchedServiceContainerError(reflect.TypeOf(service), reflect.TypeOf(key))
	}
	return repo, nil
}

func (s *ServiceContainerManager) GetUserMongoRepository() (mongodb_interface.User, error) {
	key := (*mongodb_interface.User)(nil)
	service, err := s.Get(reflect.TypeOf(key))
	if err != nil {
		return nil, errors.NewServiceWasNotFoundIntoContainerError(reflect.TypeOf(key))
	}
	repo, ok := service.Interface().(mongodb_interface.User)
	if !ok {
		return nil, errors.NewTypesMismatchedServiceContainerError(reflect.TypeOf(service), reflect.TypeOf(key))
	}
	return repo, nil
}

func (s *ServiceContainerManager) GetBlockedTokenMongoRepository() (mongodb_interface.BlockedToken, error) {
	key := (*mongodb_interface.BlockedToken)(nil)
	service, err := s.Get(reflect.TypeOf(key))
	if err != nil {
		return nil, errors.NewServiceWasNotFoundIntoContainerError(reflect.TypeOf(key))
	}
	repo, ok := service.Interface().(mongodb_interface.BlockedToken)
	if !ok {
		return nil, errors.NewTypesMismatchedServiceContainerError(reflect.TypeOf(service), reflect.TypeOf(key))
	}
	return repo, nil
}

func (s *ServiceContainerManager) GetResourceCacheRepository() (cache_interface.Resource, error) {
	key := (*cache_interface.Resource)(nil)
	service, err := s.Get(reflect.TypeOf(key))
	if err != nil {
		return nil, errors.NewServiceWasNotFoundIntoContainerError(reflect.TypeOf(key))
	}
	repo, ok := service.Interface().(cache_interface.Resource)
	if !ok {
		return nil, errors.NewTypesMismatchedServiceContainerError(reflect.TypeOf(service), reflect.TypeOf(key))
	}
	return repo, nil
}

func (s *ServiceContainerManager) GetVideoCacheRepository() (cache_interface.Video, error) {
	key := (*cache_interface.Video)(nil)
	service, err := s.Get(reflect.TypeOf(key))
	if err != nil {
		return nil, errors.NewServiceWasNotFoundIntoContainerError(reflect.TypeOf(key))
	}
	repo, ok := service.Interface().(cache_interface.Video)
	if !ok {
		return nil, errors.NewTypesMismatchedServiceContainerError(reflect.TypeOf(service), reflect.TypeOf(key))
	}
	return repo, nil
}

func (s *ServiceContainerManager) GetUserCacheRepository() (cache_interface.User, error) {
	key := (*cache_interface.User)(nil)
	service, err := s.Get(reflect.TypeOf(key))
	if err != nil {
		return nil, errors.NewServiceWasNotFoundIntoContainerError(reflect.TypeOf(key))
	}
	repo, ok := service.Interface().(cache_interface.User)
	if !ok {
		return nil, errors.NewTypesMismatchedServiceContainerError(reflect.TypeOf(service), reflect.TypeOf(key))
	}
	return repo, nil
}

func (s *ServiceContainerManager) GetAccessService() (accessor_interface.Accessor, error) {
	key := (*accessor_interface.Accessor)(nil)
	service, err := s.Get(reflect.TypeOf(key))
	if err != nil {
		return nil, errors.NewServiceWasNotFoundIntoContainerError(reflect.TypeOf(key))
	}
	accessService, ok := service.Interface().(accessor_interface.Accessor)
	if !ok {
		return nil, errors.NewTypesMismatchedServiceContainerError(reflect.TypeOf(service), reflect.TypeOf(key))
	}
	return accessService, nil
}

func (s *ServiceContainerManager) GetResourceBuilder() (builder_interface.Resource, error) {
	key := (*builder_interface.Resource)(nil)
	service, err := s.Get(reflect.TypeOf(key))
	if err != nil {
		return nil, errors.NewServiceWasNotFoundIntoContainerError(reflect.TypeOf(key))
	}
	builderService, ok := service.Interface().(builder_interface.Resource)
	if !ok {
		return nil, errors.NewTypesMismatchedServiceContainerError(reflect.TypeOf(service), reflect.TypeOf(key))
	}
	return builderService, nil
}

func (s *ServiceContainerManager) GetResourceValidator() (validator_interface.Resource, error) {
	key := (*validator_interface.Resource)(nil)
	service, err := s.Get(reflect.TypeOf(key))
	if err != nil {
		return nil, errors.NewServiceWasNotFoundIntoContainerError(reflect.TypeOf(key))
	}
	validatorService, ok := service.Interface().(validator_interface.Resource)
	if !ok {
		return nil, errors.NewTypesMismatchedServiceContainerError(reflect.TypeOf(service), reflect.TypeOf(key))
	}
	return validatorService, nil
}

func (s *ServiceContainerManager) GetResourceRepository() (repository_interface.Resource, error) {
	key := (*repository_interface.Resource)(nil)
	service, err := s.Get(reflect.TypeOf(key))
	if err != nil {
		return nil, errors.NewServiceWasNotFoundIntoContainerError(reflect.TypeOf(key))
	}
	repo, ok := service.Interface().(repository_interface.Resource)
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

func (s *ServiceContainerManager) GetBlockedTokenRepository() (repository_interface.BlockedToken, error) {
	key := (*repository_interface.BlockedToken)(nil)
	reflectService, err := s.Get(reflect.TypeOf(key))
	if err != nil {
		return nil, errors.NewServiceWasNotFoundIntoContainerError(reflect.TypeOf(key))
	}
	service, ok := reflectService.Interface().(repository_interface.BlockedToken)
	if !ok {
		return nil, errors.NewTypesMismatchedServiceContainerError(reflect.TypeOf(reflectService), reflect.TypeOf(key))
	}
	return service, nil
}

func (s *ServiceContainerManager) GetVideoBuilder() (builder_interface.Video, error) {
	key := (*builder_interface.Video)(nil)
	reflectService, err := s.Get(reflect.TypeOf(key))
	if err != nil {
		return nil, errors.NewServiceWasNotFoundIntoContainerError(reflect.TypeOf(key))
	}
	service, ok := reflectService.Interface().(builder_interface.Video)
	if !ok {
		return nil, errors.NewTypesMismatchedServiceContainerError(reflect.TypeOf(reflectService), reflect.TypeOf(key))
	}
	return service, nil
}

func (s *ServiceContainerManager) GetVideoValidator() (validator_interface.Video, error) {
	key := (*validator_interface.Video)(nil)
	reflectService, err := s.Get(reflect.TypeOf(key))
	if err != nil {
		return nil, errors.NewServiceWasNotFoundIntoContainerError(reflect.TypeOf(key))
	}
	service, ok := reflectService.Interface().(validator_interface.Video)
	if !ok {
		return nil, errors.NewTypesMismatchedServiceContainerError(reflect.TypeOf(reflectService), reflect.TypeOf(key))
	}
	return service, nil
}

func (s *ServiceContainerManager) GetVideoRepository() (repository_interface.Video, error) {
	key := (*repository_interface.Video)(nil)
	reflectService, err := s.Get(reflect.TypeOf(key))
	if err != nil {
		return nil, errors.NewServiceWasNotFoundIntoContainerError(reflect.TypeOf(key))
	}
	service, ok := reflectService.Interface().(repository_interface.Video)
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

func (s *ServiceContainerManager) GetUserBuilder() (builder_interface.User, error) {
	key := (*builder_interface.User)(nil)
	reflectService, err := s.Get(reflect.TypeOf(key))
	if err != nil {
		return nil, errors.NewServiceWasNotFoundIntoContainerError(reflect.TypeOf(key))
	}
	service, ok := reflectService.Interface().(builder_interface.User)
	if !ok {
		return nil, errors.NewTypesMismatchedServiceContainerError(reflect.TypeOf(reflectService), reflect.TypeOf(key))
	}
	return service, nil
}

func (s *ServiceContainerManager) GetUserValidator() (validator_interface.User, error) {
	key := (*validator_interface.User)(nil)
	reflectService, err := s.Get(reflect.TypeOf(key))
	if err != nil {
		return nil, errors.NewServiceWasNotFoundIntoContainerError(reflect.TypeOf(key))
	}
	service, ok := reflectService.Interface().(validator_interface.User)
	if !ok {
		return nil, errors.NewTypesMismatchedServiceContainerError(reflect.TypeOf(reflectService), reflect.TypeOf(key))
	}
	return service, nil
}

func (s *ServiceContainerManager) GetUserRepository() (repository_interface.User, error) {
	key := (*repository_interface.User)(nil)
	reflectService, err := s.Get(reflect.TypeOf(key))
	if err != nil {
		return nil, errors.NewServiceWasNotFoundIntoContainerError(reflect.TypeOf(key))
	}
	service, ok := reflectService.Interface().(repository_interface.User)
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

func (s *ServiceContainerManager) GetAuthBuilder() (builder_interface.Auth, error) {
	key := (*builder_interface.Auth)(nil)
	reflectService, err := s.Get(reflect.TypeOf(key))
	if err != nil {
		return nil, errors.NewServiceWasNotFoundIntoContainerError(reflect.TypeOf(key))
	}
	service, ok := reflectService.Interface().(builder_interface.Auth)
	if !ok {
		return nil, errors.NewTypesMismatchedServiceContainerError(reflect.TypeOf(reflectService), reflect.TypeOf(key))
	}
	return service, nil
}

func (s *ServiceContainerManager) GetAuthValidator() (validator_interface.Auth, error) {
	key := (*validator_interface.Auth)(nil)
	reflectService, err := s.Get(reflect.TypeOf(key))
	if err != nil {
		return nil, errors.NewServiceWasNotFoundIntoContainerError(reflect.TypeOf(key))
	}
	service, ok := reflectService.Interface().(validator_interface.Auth)
	if !ok {
		return nil, errors.NewTypesMismatchedServiceContainerError(reflect.TypeOf(reflectService), reflect.TypeOf(key))
	}
	return service, nil
}

func (s *ServiceContainerManager) GetAuthService() (authenticator_interface.Authenticator, error) {
	key := (*authenticator_interface.Authenticator)(nil)
	reflectService, err := s.Get(reflect.TypeOf(key))
	if err != nil {
		return nil, errors.NewServiceWasNotFoundIntoContainerError(reflect.TypeOf(key))
	}
	service, ok := reflectService.Interface().(authenticator_interface.Authenticator)
	if !ok {
		return nil, errors.NewTypesMismatchedServiceContainerError(reflect.TypeOf(reflectService), reflect.TypeOf(key))
	}
	return service, nil
}

func (s *ServiceContainerManager) GetLoggerService() (logger_interface.Logger, error) {
	key := (*logger_interface.Logger)(nil)
	reflectService, err := s.Get(reflect.TypeOf(key))
	if err != nil {
		return nil, errors.NewServiceWasNotFoundIntoContainerError(reflect.TypeOf(key))
	}
	service, ok := reflectService.Interface().(logger_interface.Logger)
	if !ok {
		return nil, errors.NewTypesMismatchedServiceContainerError(reflect.TypeOf(reflectService), reflect.TypeOf(key))
	}
	return service, nil
}

func (s *ServiceContainerManager) GetCacheService() (cacher_interface.Cacher, error) {
	key := (*cacher_interface.Cacher)(nil)
	reflectService, err := s.Get(reflect.TypeOf(key))
	if err != nil {
		return nil, errors.NewServiceWasNotFoundIntoContainerError(reflect.TypeOf(key))
	}
	service, ok := reflectService.Interface().(cacher_interface.Cacher)
	if !ok {
		return nil, errors.NewTypesMismatchedServiceContainerError(reflect.TypeOf(reflectService), reflect.TypeOf(key))
	}
	return service, nil
}

func (s *ServiceContainerManager) GetFileStorageService() (file_interface.Storage, error) {
	key := (*file_interface.Storage)(nil)
	reflectService, err := s.Get(reflect.TypeOf(key))
	if err != nil {
		return nil, errors.NewServiceWasNotFoundIntoContainerError(reflect.TypeOf(key))
	}
	service, ok := reflectService.Interface().(file_interface.Storage)
	if !ok {
		return nil, errors.NewTypesMismatchedServiceContainerError(reflect.TypeOf(reflectService), reflect.TypeOf(key))
	}
	return service, nil
}

func (s *ServiceContainerManager) GetFileNameComputerService() (file_interface.NameComputer, error) {
	key := (*file_interface.NameComputer)(nil)
	reflectService, err := s.Get(reflect.TypeOf(key))
	if err != nil {
		return nil, errors.NewServiceWasNotFoundIntoContainerError(reflect.TypeOf(key))
	}
	service, ok := reflectService.Interface().(file_interface.NameComputer)
	if !ok {
		return nil, errors.NewTypesMismatchedServiceContainerError(reflect.TypeOf(reflectService), reflect.TypeOf(key))
	}
	return service, nil
}

func (s *ServiceContainerManager) GetFileUploaderService() (uploader_interface.Uploader, error) {
	key := (*uploader_interface.Uploader)(nil)
	reflectService, err := s.Get(reflect.TypeOf(key))
	if err != nil {
		return nil, errors.NewServiceWasNotFoundIntoContainerError(reflect.TypeOf(key))
	}
	service, ok := reflectService.Interface().(uploader_interface.Uploader)
	if !ok {
		return nil, errors.NewTypesMismatchedServiceContainerError(reflect.TypeOf(reflectService), reflect.TypeOf(key))
	}
	return service, nil
}

func (s *ServiceContainerManager) GetFileReaderService() (reader_interface.FileReader, error) {
	key := (*reader_interface.FileReader)(nil)
	reflectService, err := s.Get(reflect.TypeOf(key))
	if err != nil {
		return nil, errors.NewServiceWasNotFoundIntoContainerError(reflect.TypeOf(key))
	}
	service, ok := reflectService.Interface().(reader_interface.FileReader)
	if !ok {
		return nil, errors.NewTypesMismatchedServiceContainerError(reflect.TypeOf(reflectService), reflect.TypeOf(key))
	}
	return service, nil
}

func (s *ServiceContainerManager) GetRequestParametersExtractorService() (extractor_interface.RequestParams, error) {
	key := (*extractor_interface.RequestParams)(nil)
	reflectService, err := s.Get(reflect.TypeOf(key))
	if err != nil {
		return nil, errors.NewServiceWasNotFoundIntoContainerError(reflect.TypeOf(key))
	}
	service, ok := reflectService.Interface().(extractor_interface.RequestParams)
	if !ok {
		return nil, errors.NewTypesMismatchedServiceContainerError(reflect.TypeOf(reflectService), reflect.TypeOf(key))
	}
	return service, nil
}

func (s *ServiceContainerManager) GetResponderService() (response_interface.Responder, error) {
	key := (*response_interface.Responder)(nil)
	reflectService, err := s.Get(reflect.TypeOf(key))
	if err != nil {
		return nil, errors.NewServiceWasNotFoundIntoContainerError(reflect.TypeOf(key))
	}
	service, ok := reflectService.Interface().(response_interface.Responder)
	if !ok {
		return nil, errors.NewTypesMismatchedServiceContainerError(reflect.TypeOf(reflectService), reflect.TypeOf(key))
	}
	return service, nil
}

func (s *ServiceContainerManager) GetPasswordHasherService() (security_interface.PasswordHasher, error) {
	key := (*security_interface.PasswordHasher)(nil)
	reflectService, err := s.Get(reflect.TypeOf(key))
	if err != nil {
		return nil, errors.NewServiceWasNotFoundIntoContainerError(reflect.TypeOf(key))
	}
	service, ok := reflectService.Interface().(security_interface.PasswordHasher)
	if !ok {
		return nil, errors.NewTypesMismatchedServiceContainerError(reflect.TypeOf(reflectService), reflect.TypeOf(key))
	}
	return service, nil
}

func (s *ServiceContainerManager) GetTokenizerService() (tokenizer_interface.Tokenizer, error) {
	key := (*tokenizer_interface.Tokenizer)(nil)
	reflectService, err := s.Get(reflect.TypeOf(key))
	if err != nil {
		return nil, errors.NewServiceWasNotFoundIntoContainerError(reflect.TypeOf(key))
	}
	service, ok := reflectService.Interface().(tokenizer_interface.Tokenizer)
	if !ok {
		return nil, errors.NewTypesMismatchedServiceContainerError(reflect.TypeOf(reflectService), reflect.TypeOf(key))
	}
	return service, nil
}

func (s *ServiceContainerManager) GetWebSocketCommunicatorService() (proto_interface.Communicator, error) {
	key := (*proto_interface.Communicator)(nil)
	reflectService, err := s.Get(reflect.TypeOf(key))
	if err != nil {
		return nil, errors.NewServiceWasNotFoundIntoContainerError(reflect.TypeOf(key))
	}
	service, ok := reflectService.Interface().(proto_interface.Communicator)
	if !ok {
		return nil, errors.NewTypesMismatchedServiceContainerError(reflect.TypeOf(reflectService), reflect.TypeOf(key))
	}
	return service, nil
}

func (s *ServiceContainerManager) GetWebSocketHandlerStrategies() ([]strategy_interface.ActionStrategy, error) {
	key := (*[]strategy_interface.ActionStrategy)(nil)
	reflectService, err := s.Get(reflect.TypeOf(key))
	if err != nil {
		return nil, errors.NewServiceWasNotFoundIntoContainerError(reflect.TypeOf(key))
	}
	service, ok := reflectService.Interface().([]strategy_interface.ActionStrategy)
	if !ok {
		return nil, errors.NewTypesMismatchedServiceContainerError(reflect.TypeOf(reflectService), reflect.TypeOf(key))
	}
	return service, nil
}

func (s *ServiceContainerManager) GetWebSocketListener() (listener_interface.ActionsListener, error) {
	key := (*listener_interface.ActionsListener)(nil)
	reflectService, err := s.Get(reflect.TypeOf(key))
	if err != nil {
		return nil, errors.NewServiceWasNotFoundIntoContainerError(reflect.TypeOf(key))
	}
	service, ok := reflectService.Interface().(listener_interface.ActionsListener)
	if !ok {
		return nil, errors.NewTypesMismatchedServiceContainerError(reflect.TypeOf(reflectService), reflect.TypeOf(key))
	}
	return service, nil
}

func (s *ServiceContainerManager) GetWebSocketHandler() (handler_interface.ActionsHandler, error) {
	key := (*handler_interface.ActionsHandler)(nil)
	reflectService, err := s.Get(reflect.TypeOf(key))
	if err != nil {
		return nil, errors.NewServiceWasNotFoundIntoContainerError(reflect.TypeOf(key))
	}
	service, ok := reflectService.Interface().(handler_interface.ActionsHandler)
	if !ok {
		return nil, errors.NewTypesMismatchedServiceContainerError(reflect.TypeOf(reflectService), reflect.TypeOf(key))
	}
	return service, nil
}

func (s *ServiceContainerManager) GetCodecsDetectorService() (detector_interface.Codecs, error) {
	key := (*detector_interface.Codecs)(nil)
	reflectService, err := s.Get(reflect.TypeOf(key))
	if err != nil {
		return nil, errors.NewServiceWasNotFoundIntoContainerError(reflect.TypeOf(key))
	}
	service, ok := reflectService.Interface().(detector_interface.Codecs)
	if !ok {
		return nil, errors.NewTypesMismatchedServiceContainerError(reflect.TypeOf(reflectService), reflect.TypeOf(key))
	}
	return service, nil
}

func (s *ServiceContainerManager) GetStreamingService() (streamer_interface.Streamer, error) {
	key := (*streamer_interface.Streamer)(nil)
	reflectService, err := s.Get(reflect.TypeOf(key))
	if err != nil {
		return nil, errors.NewServiceWasNotFoundIntoContainerError(reflect.TypeOf(key))
	}
	service, ok := reflectService.Interface().(streamer_interface.Streamer)
	if !ok {
		return nil, errors.NewTypesMismatchedServiceContainerError(reflect.TypeOf(reflectService), reflect.TypeOf(key))
	}
	return service, nil
}
