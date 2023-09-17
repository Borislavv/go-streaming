package service

import (
	"context"
	"github.com/Borislavv/video-streaming/internal/domain/agg"
	"github.com/Borislavv/video-streaming/internal/domain/builder"
	"github.com/Borislavv/video-streaming/internal/domain/dto"
	"github.com/Borislavv/video-streaming/internal/domain/logger"
	"github.com/Borislavv/video-streaming/internal/domain/repository"
	"github.com/Borislavv/video-streaming/internal/domain/validator"
)

type ResourceService struct {
	ctx        context.Context
	logger     logger.Logger
	uploader   Uploader
	validator  validator.Resource
	builder    builder.Resource
	repository repository.Resource
}

func NewResourceService(
	ctx context.Context,
	logger logger.Logger,
	uploader Uploader,
	validator validator.Resource,
	builder builder.Resource,
	repository repository.Resource,
) *ResourceService {
	return &ResourceService{
		ctx:        ctx,
		logger:     logger,
		uploader:   uploader,
		validator:  validator,
		builder:    builder,
		repository: repository,
	}
}

func (s *ResourceService) Upload(req dto.UploadRequest) (*agg.Resource, error) {
	if err := s.validator.ValidateUploadRequestDTO(req); err != nil {
		return nil, s.logger.LogPropagate(err)
	}

	if err := s.uploader.Upload(req); err != nil {
		return nil, s.logger.LogPropagate(err)
	}

	resource := s.builder.BuildAggFromUploadRequestDTO(req)

	if err := s.validator.ValidateAggregate(resource); err != nil {
		return nil, s.logger.LogPropagate(err)
	}

	resource, err := s.repository.Insert(s.ctx, resource)
	if err != nil {
		return nil, s.logger.LogPropagate(err)
	}

	return resource, nil
}
