package cache

import (
	"github.com/Borislavv/video-streaming/internal/domain/logger"
	"github.com/Borislavv/video-streaming/internal/domain/service/cacher"
	"github.com/Borislavv/video-streaming/internal/infrastructure/repository/storage/mongodb"
)

type VideoRepository struct {
	*mongodb.VideoRepository
	logger logger.Logger
	cache  cacher.Cacher
}

func NewVideoRepository(
	logger logger.Logger,
	cache cacher.Cacher,
	videoMongoDbRepository *mongodb.VideoRepository,
) *VideoRepository {
	return &VideoRepository{
		VideoRepository: videoMongoDbRepository,
		logger:          logger,
		cache:           cache,
	}
}
