package cacher

import "context"

type Displacer interface {
	Run(ctx context.Context, c *cache)
	Stop()
}
