package cacherinterface

type Cacher interface {
	Get(key string, fn func(CacheItem) (data interface{}, err error)) (data interface{}, err error)
	Delete(key string)
}
