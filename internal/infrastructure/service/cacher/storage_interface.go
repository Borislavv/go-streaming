package cacher

import "github.com/Borislavv/video-streaming/internal/domain/service/cacher"

type Storage interface {
	Get(key string, fn func(cacher.CacheItem) (data interface{}, err error)) (data interface{}, err error)
	Delete(key string)
	Displace()
}
