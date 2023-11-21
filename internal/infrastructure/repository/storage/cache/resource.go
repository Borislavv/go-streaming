package cache

import (
	"github.com/Borislavv/video-streaming/internal/domain/logger"
	"github.com/Borislavv/video-streaming/internal/domain/service/cacher"
	"github.com/Borislavv/video-streaming/internal/infrastructure/repository/storage/mongodb"
)

type ResourceRepository struct {
	*mongodb.ResourceRepository
	logger logger.Logger
	cache  cacher.Cacher
}

func NewResourceRepository(
	logger logger.Logger,
	cache cacher.Cacher,
	resourceMongoDbRepository *mongodb.ResourceRepository,
) *ResourceRepository {
	return &ResourceRepository{
		ResourceRepository: resourceMongoDbRepository,
		logger:             logger,
		cache:              cache,
	}
}
