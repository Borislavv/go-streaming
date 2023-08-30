package videoservice

import (
	"github.com/Borislavv/video-streaming/internal/app/service/logger"
	aggbuilder "github.com/Borislavv/video-streaming/internal/domain/builder/agg"
	dtobuilder "github.com/Borislavv/video-streaming/internal/domain/builder/dto"
	"github.com/Borislavv/video-streaming/internal/domain/repository"
)

type GetVideoService struct {
	logger          logger.Logger
	videoDtoBuilder dtobuilder.Video
	videoAggBuilder aggbuilder.Video
	videoRepository repository.Video
}

func NewGetVideoService(
	logger logger.Logger,
	videoDtoBuilder dtobuilder.Video,
	videoAggBuilder aggbuilder.Video,
	videoRepository repository.Video,
) *GetVideoService {
	return &GetVideoService{
		logger:          logger,
		videoDtoBuilder: videoDtoBuilder,
		videoAggBuilder: videoAggBuilder,
		videoRepository: videoRepository,
	}
}
