package cache

import (
	"github.com/Borislavv/video-streaming/internal/domain/logger"
	"github.com/Borislavv/video-streaming/internal/domain/service/cacher"
	"github.com/Borislavv/video-streaming/internal/infrastructure/repository/storage/mongodb"
)

type VideoRepository struct {
	*mongodb.UserRepository
	logger logger.Logger
	cache  cacher.Cacher
}
