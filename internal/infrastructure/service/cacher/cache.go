package cacher

import (
	domain_cacher_interface "github.com/Borislavv/video-streaming/internal/domain/service/cacher/interface"
	cacher_interface "github.com/Borislavv/video-streaming/internal/infrastructure/service/cacher/interface"
)

type Cache struct {
	storage   cacher_interface.Storage
	displacer cacher_interface.Displacer
}

func NewCache(storage cacher_interface.Storage, displacer cacher_interface.Displacer) *Cache {
	c := &Cache{
		storage:   storage,
		displacer: displacer,
	}
	c.displacer.Run(storage)
	return c
}

func (c *Cache) Get(key string, fn func(item domain_cacher_interface.CacheItem) (data interface{}, err error)) (data interface{}, err error) {
	return c.storage.Get(key, fn)
}

func (c *Cache) Delete(key string) {
	c.storage.Delete(key)
}
