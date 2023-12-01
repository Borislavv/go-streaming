package cache

import (
	"context"
	"fmt"
	"github.com/Borislavv/video-streaming/internal/domain/agg"
	"github.com/Borislavv/video-streaming/internal/domain/errors"
	"github.com/Borislavv/video-streaming/internal/domain/logger"
	"github.com/Borislavv/video-streaming/internal/domain/service/cacher"
	"github.com/Borislavv/video-streaming/internal/infrastructure/repository/query"
	"github.com/Borislavv/video-streaming/internal/infrastructure/repository/storage/mongodb"
	"reflect"
	"time"
)

type UserRepository struct {
	*mongodb.UserRepository
	logger logger.Logger
	cache  cacher.Cacher
}

func NewUserRepository(
	logger logger.Logger,
	cache cacher.Cacher,
	userMongoDbRepository *mongodb.UserRepository,
) *UserRepository {
	return &UserRepository{
		logger:         logger,
		cache:          cache,
		UserRepository: userMongoDbRepository,
	}
}

func (r *UserRepository) FindOneByID(ctx context.Context, q query.FindOneUserByID) (user *agg.User, err error) {
	// building a cache key
	cacheKey := fmt.Sprintf("userID_%v", q.GetID().Value.Hex())

	// fetching data from cache/storage
	userInterface, err := r.cache.Get(
		cacheKey,
		func(item cacher.CacheItem) (data interface{}, err error) {
			item.SetTTL(time.Hour)

			userAgg, err := r.UserRepository.FindOneByID(ctx, q)
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

func (r *UserRepository) FindOneByEmail(ctx context.Context, q query.FindOneUserByEmail) (user *agg.User, err error) {
	// building a cache key
	cacheKey := fmt.Sprintf("user_email_%v", q.GetEmail())

	// fetching data from cache/storage
	userInterface, err := r.cache.Get(
		cacheKey,
		func(item cacher.CacheItem) (data interface{}, err error) {
			item.SetTTL(time.Hour)

			userAgg, err := r.UserRepository.FindOneByEmail(ctx, q)
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
