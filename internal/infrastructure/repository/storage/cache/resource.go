package cache

import (
	"context"
	"encoding/json"
	"github.com/Borislavv/video-streaming/internal/domain/agg"
	"github.com/Borislavv/video-streaming/internal/domain/errors"
	"github.com/Borislavv/video-streaming/internal/domain/logger/interface"
	cacher_interface "github.com/Borislavv/video-streaming/internal/domain/service/cacher/interface"
	"github.com/Borislavv/video-streaming/internal/domain/service/di/interface"
	"github.com/Borislavv/video-streaming/internal/infrastructure/helper"
	queryinterface "github.com/Borislavv/video-streaming/internal/infrastructure/repository/query/interface"
	mongodbinterface "github.com/Borislavv/video-streaming/internal/infrastructure/repository/storage/mongodb/interface"
	"reflect"
	"time"
)

type ResourceRepository struct {
	mongodbinterface.Resource
	logger loggerinterface.Logger
	cache  cacher_interface.Cacher
}

func NewResourceRepository(serviceContainer diinterface.ContainerManager) (*ResourceRepository, error) {
	loggerService, err := serviceContainer.GetLoggerService()
	if err != nil {
		return nil, err
	}

	mongoRepository, err := serviceContainer.GetResourceMongoRepository()
	if err != nil {
		return nil, loggerService.LogPropagate(err)
	}

	cacheService, err := serviceContainer.GetCacheService()
	if err != nil {
		return nil, loggerService.LogPropagate(err)
	}

	return &ResourceRepository{
		Resource: mongoRepository,
		logger:   loggerService,
		cache:    cacheService,
	}, nil
}

func (r *ResourceRepository) FindOneByID(ctx context.Context, q queryinterface.FindOneResourceByID) (*agg.Resource, error) {
	// attempt to fetch data from cache
	if resource, err := r.findOneByID(ctx, q); err == nil {
		return resource, nil
	}
	// fetch data from storage if an error occurred
	return r.Resource.FindOneByID(ctx, q)
}

func (r *ResourceRepository) findOneByID(ctx context.Context, q queryinterface.FindOneResourceByID) (*agg.Resource, error) {
	p, err := json.Marshal(q)
	if err != nil {
		return nil, r.logger.LogPropagate(err)
	}
	cacheKey := helper.MD5(p)

	resourceInterface, err := r.cache.Get(cacheKey, func(item cacher_interface.CacheItem) (data interface{}, err error) {
		item.SetTTL(time.Hour)

		resourceAgg, err := r.Resource.FindOneByID(ctx, q)
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
