package repository

import (
	"context"
	"github.com/Borislavv/video-streaming/internal/domain/agg"
	"github.com/Borislavv/video-streaming/internal/infrastructure/repository/query"
)

type Video interface {
	FindOneByID(ctx context.Context, q query.FindOneVideoByID) (*agg.Video, error)
	FindOneByName(ctx context.Context, q query.FindOneVideoByName) (*agg.Video, error)
	FindOneByResourceID(ctx context.Context, q query.FindOneVideoByResourceID) (*agg.Video, error)
	FindList(ctx context.Context, q query.FindVideoList) (list []*agg.Video, total int64, err error)
	Insert(ctx context.Context, video *agg.Video) (*agg.Video, error)
	Update(ctx context.Context, video *agg.Video) (*agg.Video, error)
	Remove(ctx context.Context, video *agg.Video) error
}
