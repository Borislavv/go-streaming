package videoservice

import (
	"github.com/Borislavv/video-streaming/internal/app/service/logger"
	aggbuilder "github.com/Borislavv/video-streaming/internal/domain/builder/agg"
	dtobuilder "github.com/Borislavv/video-streaming/internal/domain/builder/dto"
	"github.com/Borislavv/video-streaming/internal/domain/repository"
)

type DeleteVideoService struct {
	logger          logger.Logger
	videoDtoBuilder dtobuilder.Video
	videoAggBuilder aggbuilder.Video
	videoRepository repository.Video
}

func NewDeleteVideoService(
	logger logger.Logger,
	videoDtoBuilder dtobuilder.Video,
	videoAggBuilder aggbuilder.Video,
	videoRepository repository.Video,
) *DeleteVideoService {
	return &DeleteVideoService{
		logger:          logger,
		videoDtoBuilder: videoDtoBuilder,
		videoAggBuilder: videoAggBuilder,
		videoRepository: videoRepository,
	}
}
