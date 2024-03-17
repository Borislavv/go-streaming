package cacher

import (
	"context"
	cacherinterface "github.com/Borislavv/video-streaming/internal/infrastructure/service/cacher/interface"
	"sync"
	"time"
)

type CacheDisplacer struct {
	ctx      context.Context
	cancel   context.CancelFunc
	wg       *sync.WaitGroup
	interval time.Duration
}

func NewCacheDisplacer(ctx context.Context, interval time.Duration) *CacheDisplacer {
	ctx, cancel := context.WithCancel(ctx)

	return &CacheDisplacer{
		ctx:      ctx,
		cancel:   cancel,
		wg:       &sync.WaitGroup{},
		interval: interval,
	}
}

func (d *CacheDisplacer) Run(storage cacherinterface.Storage) {
	d.wg.Add(1)
	go d.run(storage)
}

func (d *CacheDisplacer) run(storage cacherinterface.Storage) {
	ticker := time.NewTicker(d.interval)

	defer func() {
		ticker.Stop()
		d.wg.Done()
	}()

	for {
		select {
		case <-d.ctx.Done():
			go d.wg.Wait()
			return
		case <-ticker.C:
			storage.Displace()
		}
	}
}

func (d *CacheDisplacer) Stop() {
	d.cancel()
	d.wg.Wait()
}
