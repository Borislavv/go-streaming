package service

import (
	"context"
	"github.com/Borislavv/video-streaming/internal/domain/agg"
	"github.com/Borislavv/video-streaming/internal/domain/builder"
	"github.com/Borislavv/video-streaming/internal/domain/dto"
	"github.com/Borislavv/video-streaming/internal/domain/repository"
	"github.com/Borislavv/video-streaming/internal/domain/validator"
	"log"
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

func (s *VideoService) Get(req dto.GetRequest) (*agg.Video, error) {
	// validation of input request
	if err := s.validator.ValidateGetRequestDto(req); err != nil {
		return nil, err
	}

	return s.repository.Find(s.ctx, req.GetId())
}

func (s *VideoService) Create(req dto.CreateRequest) (*agg.Video, error) {
	// validation of input request
	if err := s.validator.ValidateCreateRequestDto(req); err != nil {
		return nil, err
	}

	// building an aggregate
	videoAgg := s.builder.BuildAggFromCreateRequestDto(req)

	// validation of an aggregate
	if err := s.validator.ValidateAgg(videoAgg); err != nil {
		return nil, err
	}

	// saving an aggregate into storage
	video, err := s.repository.Insert(s.ctx, videoAgg)
	if err != nil {
		s.logger.Error(err)
		return nil, err
	}

	return video, nil
}

func (s *VideoService) Update(req dto.UpdateRequest) (*agg.Video, error) {
	// validation of input request
	if err := s.validator.ValidateUpdateRequestDto(req); err != nil {
		return nil, err
	}

	// building an aggregate
	videoAgg, err := s.builder.BuildAggFromUpdateRequestDto(req)
	if err != nil {
		log.Fatalln(err)
		return nil, err
	}

	// validation of an aggregate
	if err = s.validator.ValidateAgg(videoAgg); err != nil {
		return nil, err
	}

	// saving updated aggregate into storage
	videoAgg, err = s.repository.Update(s.ctx, videoAgg)
	if err != nil {
		return nil, err
	}

	return videoAgg, nil
}
