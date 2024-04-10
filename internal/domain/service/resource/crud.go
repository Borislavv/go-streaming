package resource

import (
	"context"
	"github.com/Borislavv/video-streaming/internal/domain/agg"
	builderinterface "github.com/Borislavv/video-streaming/internal/domain/builder/interface"
	dtointerface "github.com/Borislavv/video-streaming/internal/domain/dto/interface"
	"github.com/Borislavv/video-streaming/internal/domain/logger/interface"
	repositoryinterface "github.com/Borislavv/video-streaming/internal/domain/repository/interface"
	diinterface "github.com/Borislavv/video-streaming/internal/domain/service/di/interface"
	storagerinterface "github.com/Borislavv/video-streaming/internal/domain/service/storager/interface"
	uploaderinterface "github.com/Borislavv/video-streaming/internal/domain/service/uploader/interface"
	validatorinterface "github.com/Borislavv/video-streaming/internal/domain/validator/interface"
)

type CRUDService struct {
	ctx        context.Context
	logger     loggerinterface.Logger
	uploader   uploaderinterface.Uploader
	validator  validatorinterface.Resource
	builder    builderinterface.Resource
	repository repositoryinterface.Resource
	storage    storagerinterface.Storage
}

func NewResourceService(serviceContainer diinterface.ServiceContainer) (*CRUDService, error) {
	loggerService, err := serviceContainer.GetLoggerService()
	if err != nil {
		return nil, err
	}

	ctx, err := serviceContainer.GetCtx()
	if err != nil {
		return nil, loggerService.LogPropagate(err)
	}

	uploaderService, err := serviceContainer.GetFileUploaderService()
	if err != nil {
		return nil, loggerService.LogPropagate(err)
	}

	validatorService, err := serviceContainer.GetResourceValidator()
	if err != nil {
		return nil, loggerService.LogPropagate(err)
	}

	builderService, err := serviceContainer.GetResourceBuilder()
	if err != nil {
		return nil, loggerService.LogPropagate(err)
	}

	resourceRepository, err := serviceContainer.GetResourceRepository()
	if err != nil {
		return nil, loggerService.LogPropagate(err)
	}

	fileStorageService, err := serviceContainer.GetFileStorageService()
	if err != nil {
		return nil, loggerService.LogPropagate(err)
	}

	return &CRUDService{
		ctx:        ctx,
		logger:     loggerService,
		uploader:   uploaderService,
		validator:  validatorService,
		builder:    builderService,
		repository: resourceRepository,
		storage:    fileStorageService,
	}, nil
}

// Upload - will be prepared and will be saved a file from the request. Important: the input request's DTO will
// mutate per uploading. Also, of course will be created a new instance of agg.Resource as the contract says and
// will be saved into the database.
func (s *CRUDService) Upload(req dtointerface.UploadResourceRequest) (resource *agg.Resource, err error) {
	defer func() {
		if err != nil {
			if e := s.onUploadingFailed(req); e != nil {
				s.logger.Log(e)
			}
		}
	}()

	// validation of raw uploading request
	if err = s.validator.ValidateUploadRequestDTO(req); err != nil {
		return nil, s.logger.LogPropagate(err)
	}

	// uploading the target file
	if err = s.uploader.Upload(req); err != nil {
		return nil, s.logger.LogPropagate(err)
	}

	// building resource aggregate
	resource = s.builder.BuildAggFromUploadRequestDTO(req)

	// validation of built aggregate
	if err = s.validator.ValidateAggregate(resource); err != nil {
		return nil, s.logger.LogPropagate(err)
	}

	// saving the built aggregate
	resource, err = s.repository.Insert(s.ctx, resource)
	if err != nil {
		return nil, s.logger.LogPropagate(err)
	}

	return resource, nil
}

// onUploadingFailed - will check that created file is removed.
func (s *CRUDService) onUploadingFailed(req dtointerface.UploadResourceRequest) error {
	// handle the case when the file was uploaded, but error occurred while saving an aggregate
	if req.GetUploadedFilename() != "" { // in this case, we need remove the uploaded file
		has, err := s.storage.Has(req.GetUploadedFilename())
		if err != nil {
			s.logger.Log(err)
			return s.logger.LogPropagate(err)
		}
		if has { // check that file exists, if so, then remove it
			if err = s.storage.Remove(req.GetUploadedFilename()); err != nil {
				return s.logger.LogPropagate(err)
			}
		}
	}

	return nil
}

// Delete - will remove a single video by id with dependencies.
func (s *CRUDService) Delete(req dtointerface.DeleteResourceRequest) (err error) {
	// validation of raw delete request
	if err = s.validator.ValidateDeleteRequestDTO(req); err != nil {
		return s.logger.LogPropagate(err)
	}

	// fetching the target resource aggregate
	resourceAgg, err := s.repository.FindOneByID(s.ctx, req)
	if err != nil {
		return s.logger.LogPropagate(err)
	}

	// removing the file first
	if err = s.storage.Remove(resourceAgg.Filename); err != nil {
		return s.logger.LogPropagate(err)
	}

	// removing the resource
	if err = s.repository.Remove(s.ctx, resourceAgg); err != nil {
		return s.logger.LogPropagate(err)
	}

	return nil
}
