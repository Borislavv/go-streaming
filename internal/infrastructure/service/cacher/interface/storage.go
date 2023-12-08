package cacher_interface

import cacher_interface "github.com/Borislavv/video-streaming/internal/domain/service/cacher/interface"

type Storage interface {
	Get(key string, fn func(cacher_interface.CacheItem) (data interface{}, err error)) (data interface{}, err error)
	Delete(key string)
	Displace()
}
