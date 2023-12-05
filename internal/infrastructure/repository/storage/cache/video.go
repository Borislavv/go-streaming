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
	// attempt to fetch data from cache
	if video, err := r.findOneByID(ctx, q); err == nil {
		return video, nil
	}
	// fetch data from storage if an error occurred
	return r.VideoRepository.FindOneByID(ctx, q)
}

func (r *VideoRepository) findOneByID(ctx context.Context, q query.FindOneVideoByID) (*agg.Video, error) {
	p, err := json.Marshal(q)
	if err != nil {
		return nil, r.logger.LogPropagate(err)
	}
	cacheKey := helper.MD5(p)

	// fetching data from cache/storage
	videoInterface, err := r.cache.Get(
		cacheKey,
		func(item cacher.CacheItem) (data interface{}, err error) {
			item.SetTTL(time.Hour)

			videoAgg, err := r.VideoRepository.FindOneByID(ctx, q)
			if err != nil {
				return nil, r.logger.LogPropagate(err)
			}
			return videoAgg, nil
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

func (r *VideoRepository) FindList(ctx context.Context, q query.FindVideoList) (list []*agg.Video, total int64, err error) {
	// attempt to fetch data from cache
	if list, total, err = r.findList(ctx, q); err == nil {
		return list, total, nil
	}
	// fetch data from storage if an error occurred
	return r.VideoRepository.FindList(ctx, q)
}

func (r *VideoRepository) findList(ctx context.Context, q query.FindVideoList) (list []*agg.Video, total int64, err error) {
	p, err := json.Marshal(q)
	if err != nil {
		return nil, 0, r.logger.LogPropagate(err)
	}
	cacheKey := helper.MD5(p)

	type response struct {
		List  []*agg.Video
		Total int64
	}

	responseInterface, err := r.cache.Get(
		cacheKey,
		func(item cacher.CacheItem) (data interface{}, err error) {
			item.SetTTL(time.Hour)

			l, t, e := r.VideoRepository.FindList(ctx, q)
			if err != nil {
				return nil, r.logger.LogPropagate(e)
			}

			return response{List: l, Total: t}, nil
		},
	)
	if err != nil {
		return nil, 0, r.logger.LogPropagate(err)
	}

	listResponse, ok := responseInterface.(response)
	if !ok {
		return nil, 0, errors.NewCachedDataTypeWasNotMatchedError(
			cacheKey, reflect.TypeOf(response{}), reflect.TypeOf(responseInterface),
		)
	}

	return listResponse.List, listResponse.Total, nil
}

func (r *VideoRepository) FindOneByName(ctx context.Context, q query.FindOneVideoByName) (*agg.Video, error) {
	// attempt to fetch data from cache
	if video, err := r.findOneByName(ctx, q); err == nil {
		return video, nil
	}
	// fetch data from storage if an error occurred
	return r.VideoRepository.FindOneByName(ctx, q)
}

func (r *VideoRepository) findOneByName(ctx context.Context, q query.FindOneVideoByName) (*agg.Video, error) {
	p, err := json.Marshal(q)
	if err != nil {
		return nil, r.logger.LogPropagate(err)
	}
	cacheKey := helper.MD5(p)

	videoInterface, err := r.cache.Get(cacheKey, func(item cacher.CacheItem) (data interface{}, err error) {
		item.SetTTL(time.Hour)

		videoAgg, err := r.VideoRepository.FindOneByName(ctx, q)
		if err != nil {
			return nil, r.logger.LogPropagate(err)
		}

		return videoAgg, nil
	})
	if err != nil {
		return nil, r.logger.LogPropagate(err)
	}

	videoAgg, ok := videoInterface.(*agg.Video)
	if !ok {
		return nil, errors.NewCachedDataTypeWasNotMatchedError(
			cacheKey, reflect.TypeOf(&agg.Video{}), reflect.TypeOf(videoInterface),
		)
	}

	return videoAgg, nil
}

func (r *VideoRepository) FindOneByResourceID(ctx context.Context, q query.FindOneVideoByResourceID) (*agg.Video, error) {
	// attempt to fetch data from cache
	if video, err := r.findOneByResourceID(ctx, q); err == nil {
		return video, nil
	}
	// fetch data from storage if an error occurred
	return r.VideoRepository.FindOneByResourceID(ctx, q)
}

func (r *VideoRepository) findOneByResourceID(ctx context.Context, q query.FindOneVideoByResourceID) (*agg.Video, error) {
	p, err := json.Marshal(q)
	if err != nil {
		return nil, r.logger.LogPropagate(err)
	}
	cacheKey := helper.MD5(p)

	videoInterface, err := r.cache.Get(cacheKey, func(item cacher.CacheItem) (data interface{}, err error) {
		item.SetTTL(time.Hour)

		videoAgg, err := r.VideoRepository.FindOneByResourceID(ctx, q)
		if err != nil {
			return nil, r.logger.LogPropagate(err)
		}

		return videoAgg, nil
	})
	if err != nil {
		return nil, r.logger.LogPropagate(err)
	}

	videoAgg, ok := videoInterface.(*agg.Video)
	if !ok {
		return nil, errors.NewCachedDataTypeWasNotMatchedError(
			cacheKey, reflect.TypeOf(&agg.Video{}), reflect.TypeOf(videoInterface),
		)
	}

	return videoAgg, nil
}
