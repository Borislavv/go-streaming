package cacher

import "time"

type CacheItem struct {
	data      interface{}
	addedAt   time.Time
	expiresAt time.Time
}

func NewCacheItem() *CacheItem {
	return &CacheItem{
		data:      struct{}{},
		addedAt:   time.Now(),
		expiresAt: time.Time{},
	}
}

func (i *CacheItem) SetTTL(ttl time.Duration) {
	i.expiresAt = time.Now().Add(ttl)
}

type Cache struct {
}
