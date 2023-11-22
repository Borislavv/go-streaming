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

type VideoRepository struct {
	*mongodb.VideoRepository
	logger logger.Logger
	cache  cacher.Cacher
}

func NewVideoRepository(
	logger logger.Logger,
	cache cacher.Cacher,
	videoMongoDbRepository *mongodb.VideoRepository,
) *VideoRepository {
	return &VideoRepository{
		VideoRepository: videoMongoDbRepository,
		logger:          logger,
		cache:           cache,
	}
}

func (r *VideoRepository) FindOneByID(ctx context.Context, q query.FindOneVideoByID) (*agg.Video, error) {
	// building a cache key
	cacheKey := fmt.Sprintf("videoID_%v_userID_%v ", q.GetID().Value.Hex(), q.GetUserID().Value.Hex())

	// fetching data from cache/storage
	videoInterface, err := r.cache.Get(
		cacheKey,
		func(item cacher.CacheItem) (data interface{}, err error) {
			item.SetTTL(time.Hour)

			userAgg, err := r.VideoRepository.FindOneByID(ctx, q)
			if err != nil {
				return false, r.logger.LogPropagate(err)
			}
			return userAgg, nil
		})
	if err != nil {
		return nil, r.logger.LogPropagate(err)
	}

	// casting found data to struct
	videoAgg, ok := videoInterface.(*agg.Video)
	if !ok {
		return nil, errors.NewCachedDataTypeWasNotMatchedError(
			cacheKey, reflect.TypeOf(&agg.User{}), reflect.TypeOf(videoInterface),
		)
	}

	return videoAgg, nil
}
