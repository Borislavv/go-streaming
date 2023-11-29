package cacher

import "time"

type CacheItem interface {
	SetTTL(ttl time.Duration)
}
