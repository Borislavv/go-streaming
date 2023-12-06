package cacher

type Storage interface {
	Get(key string, fn func(cacher_interface.CacheItem) (data interface{}, err error)) (data interface{}, err error)
	Delete(key string)
	Displace()
}
