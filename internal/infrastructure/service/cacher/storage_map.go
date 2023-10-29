package cacher

import (
	"context"
	"sync"
	"time"
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

func (c *MapCacheStorage) Get(key string, fn func(CacheItem) (data interface{}, err error)) (data interface{}, err error) {
	item, found := c.get(key)
	if found {
		return item.data, nil
	}

	item, err = c.compute(fn)
	if err != nil {
		return nil, err
	}

	return c.set(key, item), nil
}

func (c *MapCacheStorage) get(key string) (item *Item, found bool) {
	defer c.mu.RUnlock()
	c.mu.RLock()
	cacheItem, found := c.storage[key]
	if found {
		foundItem, ok := cacheItem.(*Item)
		if ok {
			return foundItem, true
		}
	}
	return nil, false
}

func (c *MapCacheStorage) compute(fn func(CacheItem) (data interface{}, err error)) (item *Item, err error) {
	item = NewCacheItem()
	data, err := fn(item)
	if err != nil {
		return nil, err
	}
	item.data = data
	item.addedAt = time.Now()
	return item, nil
}

func (c *MapCacheStorage) set(key string, item *Item) (data interface{}) {
	defer c.mu.Unlock()
	c.mu.Lock()
	cacheItem, found := c.storage[key]
	if found {
		foundItem, ok := cacheItem.(*Item)
		if ok {
			return foundItem.data
		}
	}
	c.storage[key] = item
	return item.data
}

func (c *MapCacheStorage) Del(key string) {
	defer c.mu.Unlock()
	c.mu.Lock()
	delete(c.storage, key)
}
