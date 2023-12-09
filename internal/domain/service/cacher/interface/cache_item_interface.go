package cacher_interface

import "time"

type CacheItem interface {
	SetTTL(ttl time.Duration)
}
