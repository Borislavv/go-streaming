package cache

import (
	"context"
	"encoding/json"
	"github.com/Borislavv/video-streaming/internal/domain/agg"
	"github.com/Borislavv/video-streaming/internal/domain/errors"
	"github.com/Borislavv/video-streaming/internal/domain/logger/interface"
	"github.com/Borislavv/video-streaming/internal/domain/service/cacher/interface"
	di_interface "github.com/Borislavv/video-streaming/internal/domain/service/di/interface"
	"github.com/Borislavv/video-streaming/internal/infrastructure/helper"
	"github.com/Borislavv/video-streaming/internal/infrastructure/repository/query/interface"
	mongodb_interface "github.com/Borislavv/video-streaming/internal/infrastructure/repository/storage/mongodb/interface"
	"reflect"
	"time"
)

type UserRepository struct {
	mongodb_interface.User
	logger logger_interface.Logger
	cache  cacher_interface.Cacher
}

func NewUserRepository(serviceContainer di_interface.ContainerManager) (*UserRepository, error) {
	loggerService, err := serviceContainer.GetLoggerService()
	if err != nil {
		return nil, err
	}

	userMongoDbRepository, err := serviceContainer.GetUserMongoRepository()
	if err != nil {
		return nil, loggerService.LogPropagate(err)
	}

	cacheService, err := serviceContainer.GetCacheService()
	if err != nil {
		return nil, loggerService.LogPropagate(err)
	}

	return &UserRepository{
		logger: loggerService,
		cache:  cacheService,
		User:   userMongoDbRepository,
	}, nil
}

func (r *UserRepository) FindOneByID(ctx context.Context, q query_interface.FindOneUserByID) (*agg.User, error) {
	// attempt to fetch data from cache
	if user, err := r.findOneByID(ctx, q); err == nil {
		return user, nil
	}
	// fetch data from storage if an error occurred
	return r.User.FindOneByID(ctx, q)
}

func (r *UserRepository) findOneByID(ctx context.Context, q query_interface.FindOneUserByID) (*agg.User, error) {
	// building a cache key
	p, err := json.Marshal(q)
	if err != nil {
		return nil, r.logger.LogPropagate(err)
	}
	cacheKey := helper.MD5(p)

	// fetching data from cache/storage
	userInterface, err := r.cache.Get(
		cacheKey,
		func(item cacher_interface.CacheItem) (data interface{}, err error) {
			item.SetTTL(time.Hour)

			userAgg, err := r.User.FindOneByID(ctx, q)
			if err != nil {
				return false, r.logger.LogPropagate(err)
			}
			return userAgg, nil
		})
	if err != nil {
		return nil, r.logger.LogPropagate(err)
	}

	// casting found data to struct
	userAgg, ok := userInterface.(*agg.User)
	if !ok {
		return nil, errors.NewCachedDataTypeWasNotMatchedError(
			cacheKey, reflect.TypeOf(&agg.User{}), reflect.TypeOf(userInterface),
		)
	}

	return userAgg, nil
}

func (r *UserRepository) FindOneByEmail(ctx context.Context, q query_interface.FindOneUserByEmail) (*agg.User, error) {
	// attempt to fetch data from cache
	if user, err := r.findOneByEmail(ctx, q); err == nil {
		return user, nil
	}
	// fetch data from storage if an error occurred
	return r.User.FindOneByEmail(ctx, q)
}

func (r *UserRepository) findOneByEmail(ctx context.Context, q query_interface.FindOneUserByEmail) (user *agg.User, err error) {
	// building a cache key
	p, err := json.Marshal(q)
	if err != nil {
		return nil, r.logger.LogPropagate(err)
	}
	cacheKey := helper.MD5(p)

	// fetching data from cache/storage
	userInterface, err := r.cache.Get(
		cacheKey,
		func(item cacher_interface.CacheItem) (data interface{}, err error) {
			item.SetTTL(time.Hour)

			userAgg, err := r.User.FindOneByEmail(ctx, q)
			if err != nil {
				return nil, r.logger.LogPropagate(err)
			}
			return userAgg, nil
		})
	if err != nil {
		return nil, r.logger.LogPropagate(err)
	}

	// casting found data to struct
	userAgg, ok := userInterface.(*agg.User)
	if !ok {
		return nil, errors.NewCachedDataTypeWasNotMatchedError(
			cacheKey, reflect.TypeOf(&agg.User{}), reflect.TypeOf(userInterface),
		)
	}

	return userAgg, nil
}
