package videoservice

import (
	"github.com/Borislavv/video-streaming/internal/app/service/logger"
	aggbuilder "github.com/Borislavv/video-streaming/internal/domain/builder/agg"
	dtobuilder "github.com/Borislavv/video-streaming/internal/domain/builder/dto"
	"github.com/Borislavv/video-streaming/internal/domain/repository"
)

type ListVideoService struct {
	logger          logger.Logger
	videoDtoBuilder dtobuilder.Video
	videoAggBuilder aggbuilder.Video
	videoRepository repository.Video
}

func NewListVideoService(
	logger logger.Logger,
	videoDtoBuilder dtobuilder.Video,
	videoAggBuilder aggbuilder.Video,
	videoRepository repository.Video,
) *ListVideoService {
	return &ListVideoService{
		logger:          logger,
		videoDtoBuilder: videoDtoBuilder,
		videoAggBuilder: videoAggBuilder,
		videoRepository: videoRepository,
	}
}
