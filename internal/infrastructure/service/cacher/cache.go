package cacher

import (
	domain_cacherinterface "github.com/Borislavv/video-streaming/internal/domain/service/cacher/interface"
	cacherinterface "github.com/Borislavv/video-streaming/internal/infrastructure/service/cacher/interface"
)

type Cache struct {
	storage   cacherinterface.Storage
	displacer cacherinterface.Displacer
}

func NewCache(storage cacherinterface.Storage, displacer cacherinterface.Displacer) *Cache {
	c := &Cache{
		storage:   storage,
		displacer: displacer,
	}
	c.displacer.Run(storage)
	return c
}

func (c *Cache) Get(key string, fn func(item domain_cacherinterface.CacheItem) (data interface{}, err error)) (data interface{}, err error) {
	return c.storage.Get(key, fn)
}

func (c *Cache) Delete(key string) {
	c.storage.Delete(key)
}
