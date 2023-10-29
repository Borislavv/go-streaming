package cacher

type Storage interface {
	Get(key string, callable func(item CacheItem) (data interface{}, err error))
	Del(key string)
}
