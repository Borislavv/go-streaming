package resource

import (
	"context"
	"github.com/Borislavv/video-streaming/internal/domain/agg"
	"github.com/Borislavv/video-streaming/internal/domain/builder"
	"github.com/Borislavv/video-streaming/internal/domain/dto"
	"github.com/Borislavv/video-streaming/internal/domain/logger"
	"github.com/Borislavv/video-streaming/internal/domain/repository"
	"github.com/Borislavv/video-streaming/internal/domain/service/storager"
	"github.com/Borislavv/video-streaming/internal/domain/service/uploader"
	"github.com/Borislavv/video-streaming/internal/domain/validator"
)

type CRUDService struct {
	ctx        context.Context
	logger     logger.Logger
	uploader   uploader.Uploader
	validator  validator.Resource
	builder    builder.Resource
	repository repository.Resource
	storage    storager.Storage
}

func NewResourceService(
	ctx context.Context,
	logger logger.Logger,
	uploader uploader.Uploader,
	validator validator.Resource,
	builder builder.Resource,
	repository repository.Resource,
	storage storager.Storage,
) *CRUDService {
	return &CRUDService{
		ctx:        ctx,
		logger:     logger,
		uploader:   uploader,
		validator:  validator,
		builder:    builder,
		repository: repository,
		storage:    storage,
	}
}

// Upload - will be prepared and  will be saved a file from the request. Important: the input request's DTO will
// mutate per uploading. Also, of course will be created a new instance of agg.Resource as the contract says and
// will be saved into the database.
func (s *CRUDService) Upload(req dto.UploadResourceRequest) (resource *agg.Resource, err error) {
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
func (s *CRUDService) onUploadingFailed(req dto.UploadResourceRequest) error {
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
func (s *CRUDService) Delete(req dto.DeleteResourceRequest) (err error) {
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
