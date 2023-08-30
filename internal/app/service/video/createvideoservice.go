package videoservice

import (
	"context"
	"github.com/Borislavv/video-streaming/internal/app/service/logger"
	aggbuilder "github.com/Borislavv/video-streaming/internal/domain/builder/agg"
	"github.com/Borislavv/video-streaming/internal/domain/dto"
	"github.com/Borislavv/video-streaming/internal/domain/repository"
	aggvalidator "github.com/Borislavv/video-streaming/internal/domain/validator/agg"
	dtovalidator "github.com/Borislavv/video-streaming/internal/domain/validator/dto"
)

type CreateVideoService struct {
	ctx               context.Context
	logger            logger.Logger
	videoDtoValidator dtovalidator.Video
	videoAggValidator aggvalidator.Video
	videoAggBuilder   aggbuilder.Video
	videoRepository   repository.Video
}

func NewCreateVideoService(
	ctx context.Context,
	logger logger.Logger,
	videoDtoValidator dtovalidator.Video,
	videoAggValidator aggvalidator.Video,
	videoAggBuilder aggbuilder.Video,
	videoRepository repository.Video,
) *CreateVideoService {
	return &CreateVideoService{
		ctx:               ctx,
		logger:            logger,
		videoDtoValidator: videoDtoValidator,
		videoAggValidator: videoAggValidator,
		videoAggBuilder:   videoAggBuilder,
		videoRepository:   videoRepository,
	}
}

func (s *CreateVideoService) Create(video *dto.Video) (string, error) {
	// validation of input request
	if err := s.videoDtoValidator.Validate(video); err != nil {
		return "", err
	}

	// building an aggregate
	agg := s.videoAggBuilder.Build(video)

	// validation of aggregate
	if err := s.videoAggValidator.Validate(agg); err != nil {
		return "", err
	}

	// saving an aggregate into storage
	id, err := s.videoRepository.Insert(s.ctx, agg)
	if err != nil {
		s.logger.Error(err)
		return "", err
	}

	return id, nil
}
