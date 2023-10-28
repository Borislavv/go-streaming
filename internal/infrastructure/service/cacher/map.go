package cacher

import (
	"github.com/Borislavv/video-streaming/internal/domain/logger"
	"sync"
)

type Item struct {
	Data       interface{}
	Expiration int64
}

type cache struct {
}

type MapCache struct {
	*cache

	mu      sync.Mutex
	logger  logger.Logger
	storage map[string]interface{}
}
