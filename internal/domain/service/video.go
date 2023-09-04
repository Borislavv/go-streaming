package service

import (
	"context"
	"github.com/Borislavv/video-streaming/internal/domain/builder"
	"github.com/Borislavv/video-streaming/internal/domain/dto"
	"github.com/Borislavv/video-streaming/internal/domain/repository"
	"github.com/Borislavv/video-streaming/internal/domain/validator"
	"github.com/Borislavv/video-streaming/internal/domain/vo"
)

type VideoService struct {
	ctx        context.Context
	logger     Logger
	builder    builder.Video
	validator  validator.Video
	repository repository.Video
}

func NewVideoService(
	ctx context.Context,
	logger Logger,
	builder builder.Video,
	validator validator.Video,
	repository repository.Video,
) *VideoService {
	return &VideoService{
		ctx:        ctx,
		logger:     logger,
		builder:    builder,
		validator:  validator,
		repository: repository,
	}
}

func (s *VideoService) Create(video dto.CreateRequest) (*vo.ID, error) {
	// validation of input request
	if err := s.validator.ValidateCreateRequestDto(video); err != nil {
		return nil, err
	}

	// building an aggregate
	agg := s.builder.BuildAggFromCreateRequestDto(video)

	// validation of aggregate
	if err := s.validator.ValidateAgg(agg); err != nil {
		return nil, err
	}

	// saving an aggregate into storage
	id, err := s.repository.Insert(s.ctx, agg)
	if err != nil {
		s.logger.Error(err)
		return nil, err
	}

	return id, nil
}
