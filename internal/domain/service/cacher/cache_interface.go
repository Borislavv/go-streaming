package cacher

import "github.com/Borislavv/video-streaming/internal/infrastructure/service/cacher"

type Cacher interface {
	Get(key string, fn func(item cacher.CacheItem) (data interface{}, err error)) (data interface{}, err error)
	Delete(key string)
}
