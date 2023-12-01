package cache

import (
	"context"
	"encoding/json"
	"github.com/Borislavv/video-streaming/internal/domain/agg"
	"github.com/Borislavv/video-streaming/internal/domain/errors"
	"github.com/Borislavv/video-streaming/internal/domain/logger"
	"github.com/Borislavv/video-streaming/internal/domain/service/cacher"
	"github.com/Borislavv/video-streaming/internal/infrastructure/helper"
	"github.com/Borislavv/video-streaming/internal/infrastructure/repository/query"
	"github.com/Borislavv/video-streaming/internal/infrastructure/repository/storage/mongodb"
	"reflect"
	"time"
)

type ResourceRepository struct {
	*mongodb.ResourceRepository
	logger logger.Logger
	cache  cacher.Cacher
}

func NewResourceRepository(
	logger logger.Logger,
	cache cacher.Cacher,
	resourceMongoDbRepository *mongodb.ResourceRepository,
) *ResourceRepository {
	return &ResourceRepository{
		ResourceRepository: resourceMongoDbRepository,
		logger:             logger,
		cache:              cache,
	}
}

func (r *ResourceRepository) FindOneByID(ctx context.Context, q query.FindOneResourceByID) (*agg.Resource, error) {
	// attempt to fetch data from cache
	if resource, err := r.findOneByID(ctx, q); err == nil {
		return resource, nil
	}
	// fetch data from storage if an error occurred
	return r.ResourceRepository.FindOneByID(ctx, q)
}

func (r *ResourceRepository) findOneByID(ctx context.Context, q query.FindOneResourceByID) (*agg.Resource, error) {
	p, err := json.Marshal(q)
	if err != nil {
		return nil, r.logger.LogPropagate(err)
	}
	cacheKey := helper.MD5(p)

	resourceInterface, err := r.cache.Get(cacheKey, func(item cacher.CacheItem) (data interface{}, err error) {
		item.SetTTL(time.Hour)

		resourceAgg, err := r.ResourceRepository.FindOneByID(ctx, q)
		if err != nil {
			return nil, r.logger.LogPropagate(err)
		}

		return resourceAgg, nil
	})
	if err != nil {
		return nil, r.logger.LogPropagate(err)
	}

	resourceAgg, ok := resourceInterface.(*agg.Resource)
	if !ok {
		return nil, errors.NewCachedDataTypeWasNotMatchedError(
			cacheKey, reflect.TypeOf(&agg.Resource{}), reflect.TypeOf(resourceInterface),
		)
	}

	return resourceAgg, nil
}
