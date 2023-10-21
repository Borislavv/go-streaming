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

func (s *CRUDService) Upload(reqDTO dto.UploadResourceRequest) (resource *agg.Resource, err error) {
	if err = s.validator.ValidateUploadRequestDTO(reqDTO); err != nil {
		return nil, s.logger.LogPropagate(err)
	}

	if err = s.uploader.Upload(reqDTO); err != nil {
		return nil, s.logger.LogPropagate(err)
	}

	resource = s.builder.BuildAggFromUploadRequestDTO(reqDTO)

	if err = s.validator.ValidateAggregate(resource); err != nil {
		return nil, s.logger.LogPropagate(err)
	}

	resource, err = s.repository.Insert(s.ctx, resource)
	if err != nil {
		return nil, s.logger.LogPropagate(err)
	}

	return resource, nil
}

func (s *CRUDService) Delete(reqDTO dto.DeleteResourceRequest) (err error) {
	// validation of raw delete request
	if err = s.validator.ValidateDeleteRequestDTO(reqDTO); err != nil {
		return s.logger.LogPropagate(err)
	}

	// fetching the target resource aggregate
	resourceAgg, err := s.repository.Find(s.ctx, reqDTO.GetID())
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
