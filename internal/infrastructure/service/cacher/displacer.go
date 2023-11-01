package cacher

import (
	"context"
	"github.com/Borislavv/video-streaming/internal/domain/logger"
	"sync"
	"time"
)

type CacheDisplacer struct {
	logger   logger.Logger
	ctx      context.Context
	once     *sync.Once
	wg       *sync.WaitGroup
	interval time.Duration
	cancel   context.CancelFunc
}

func NewCacheDisplacer(logger logger.Logger, ctx context.Context, interval time.Duration) *CacheDisplacer {
	ctx, cancel := context.WithCancel(ctx)

	return &CacheDisplacer{
		logger:   logger,
		ctx:      ctx,
		cancel:   cancel,
		once:     &sync.Once{},
		wg:       &sync.WaitGroup{},
		interval: interval,
	}
}

func (d *CacheDisplacer) Run(storage Storage) {
	d.wg.Add(1)
	go d.run(storage)
}

func (d *CacheDisplacer) run(storage Storage) {
	ticker := time.NewTicker(d.interval)

	defer func() {
		ticker.Stop()
		d.wg.Done()
	}()

	for {
		select {
		case <-d.ctx.Done():
			d.logger.Info("RUN FINISHED BY DONE CH")
			go func() {
				d.wg.Wait()
				d.logger.Info("FULL STOPPED IN ANOTHER GOROUTINE")
			}()
			return
		case <-ticker.C:
			storage.Displace()
		}
	}
}

func (d *CacheDisplacer) Stop() {
	// broadcasting `stop` action by closing chan
	d.cancel()
	d.wg.Wait()
	d.logger.Info("FULL STOPPED IN STOP METHOD")
}
