package cacher

import (
	"context"
	"sync"
)

type MapCacheStorage struct {
	*mapCacheStorage
}
type mapCacheStorage struct {
	ctx     context.Context
	mu      sync.RWMutex
	storage map[string]CacheItem
}

func NewMapCacheStorage(ctx context.Context) *MapCacheStorage {
	return &MapCacheStorage{
		mapCacheStorage: &mapCacheStorage{
			ctx:     ctx,
			mu:      sync.RWMutex{},
			storage: map[string]CacheItem{},
		},
	}
}

func (c *MapCacheStorage) Get(key string, callable func() (data interface{}, err error)) {

}