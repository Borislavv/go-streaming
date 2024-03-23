package cacher

import (
	"context"
	cacherinterface "github.com/Borislavv/video-streaming/internal/domain/service/cacher/interface"
	"sync"
	"time"
)

type MapCacheStorage struct {
	*mapCacheStorage
}
type mapCacheStorage struct {
	ctx      context.Context
	mu       sync.RWMutex
	storage  map[string]*Item
	capacity int64
}

// NewMapCacheStorage is a constructor of MapCacheStorage structure.
func NewMapCacheStorage(ctx context.Context) *MapCacheStorage {
	return &MapCacheStorage{
		mapCacheStorage: &mapCacheStorage{
			ctx:     ctx,
			mu:      sync.RWMutex{},
			storage: map[string]*Item{},
		},
	}
}

func (c *MapCacheStorage) Get(key string, fn func(cacherinterface.CacheItem) (data interface{}, err error)) (data interface{}, err error) {
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
	item, found = c.storage[key]
	return item, found
}

func (c *MapCacheStorage) compute(fn func(cacherinterface.CacheItem) (data interface{}, err error)) (item *Item, err error) {
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
		return cacheItem.data
	}
	c.storage[key] = item
	return item.data
}

func (c *MapCacheStorage) Delete(key string) {
	defer c.mu.Unlock()
	c.mu.Lock()
	delete(c.storage, key)
}

func (c *MapCacheStorage) Displace() {
	var keys []string

	c.mu.RLock()
	for key, item := range c.storage {
		if !item.expiresAt.IsZero() && item.expiresAt.UnixNano() <= time.Now().UnixNano() {
			keys = append(keys, key)
		}
	}
	c.mu.RUnlock()

	c.mu.Lock()
	for _, key := range keys {
		delete(c.storage, key)
	}
	c.mu.Unlock()
}
