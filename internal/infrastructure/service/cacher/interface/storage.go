package cacherinterface

import cacherinterface "github.com/Borislavv/video-streaming/internal/domain/service/cacher/interface"

type Storage interface {
	Get(key string, fn func(cacherinterface.CacheItem) (data interface{}, err error)) (data interface{}, err error)
	Delete(key string)
	Displace()
}
