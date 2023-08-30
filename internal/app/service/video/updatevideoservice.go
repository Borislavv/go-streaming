package videoservice

import (
	"github.com/Borislavv/video-streaming/internal/app/service/logger"
	aggbuilder "github.com/Borislavv/video-streaming/internal/domain/builder/agg"
	dtobuilder "github.com/Borislavv/video-streaming/internal/domain/builder/dto"
	"github.com/Borislavv/video-streaming/internal/domain/repository"
)

type UpdateVideoService struct {
	logger          logger.Logger
	videoDtoBuilder dtobuilder.Video
	videoAggBuilder aggbuilder.Video
	videoRepository repository.Video
}

func NewUpdateVideoService(
	logger logger.Logger,
	videoDtoBuilder dtobuilder.Video,
	videoAggBuilder aggbuilder.Video,
	videoRepository repository.Video,
) *UpdateVideoService {
	return &UpdateVideoService{
		logger:          logger,
		videoDtoBuilder: videoDtoBuilder,
		videoAggBuilder: videoAggBuilder,
		videoRepository: videoRepository,
	}
}
