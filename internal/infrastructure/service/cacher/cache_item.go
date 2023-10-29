package cacher

import "time"

type Item struct {
	data      interface{}
	addedAt   time.Time
	expiresAt time.Time
}

func NewCacheItem() *Item {
	return &Item{
		data:      struct{}{},
		addedAt:   time.Now(),
		expiresAt: time.Time{},
	}
}

func (i *Item) SetTTL(ttl time.Duration) {
	i.expiresAt = time.Now().Add(ttl)
}
