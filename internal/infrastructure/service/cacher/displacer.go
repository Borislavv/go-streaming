package cacher

import (
	"context"
	"sync"
	"time"
)

type CacheDisplacer struct {
	ctx      context.Context
	once     *sync.Once
	wg       *sync.WaitGroup
	interval time.Duration
	stopCh   chan struct{}
	doneCh   chan struct{}
}

func NewCacheDisplacer(ctx context.Context, interval time.Duration) *CacheDisplacer {
	return &CacheDisplacer{
		ctx:      ctx,
		once:     &sync.Once{},
		wg:       &sync.WaitGroup{},
		interval: interval,
		stopCh:   make(chan struct{}, 1),
		doneCh:   make(chan struct{}, 1),
	}
}

func (d *CacheDisplacer) Run(storage Storage) {
	d.wg.Add(2)
	go d.run(storage)
	go d.listenStop()
}

func (d *CacheDisplacer) run(storage Storage) {
	ticker := time.NewTicker(d.interval)

	defer func() {
		ticker.Stop()
		d.wg.Done()
	}()

	for {
		select {
		case <-d.doneCh:
			return
		case <-ticker.C:
			storage.Displace()
		}
	}
}

func (d *CacheDisplacer) Stop() {
	// broadcasting `stop` action by closing chan
	close(d.stopCh)
	d.wg.Wait()
}

// stop is a simple fun-in pattern.
func (d *CacheDisplacer) listenStop() {
	defer func() {
		// broadcasting `done` action by closing chan
		close(d.doneCh)
		d.wg.Done()
	}()

	for {
		select {
		case <-d.ctx.Done():
			// awaiting all goroutines will be stopped in another goroutine
			go d.Stop()
			return
		case <-d.stopCh:
			// used Stop func., WaitGroup will be awaited by Stop func.
			return
		}
	}
}
