package cache_interface

import (
	"context"
	"github.com/Borislavv/video-streaming/internal/domain/agg"
	query_interface "github.com/Borislavv/video-streaming/internal/infrastructure/repository/query/interface"
)

type Video interface {
	FindOneByID(ctx context.Context, q query_interface.FindOneVideoByID) (*agg.Video, error)
	FindOneByName(ctx context.Context, q query_interface.FindOneVideoByName) (*agg.Video, error)
	FindOneByResourceID(ctx context.Context, q query_interface.FindOneVideoByResourceID) (*agg.Video, error)
	FindList(ctx context.Context, q query_interface.FindVideoList) (list []*agg.Video, total int64, err error)
	Insert(ctx context.Context, video *agg.Video) (*agg.Video, error)
	Update(ctx context.Context, video *agg.Video) (*agg.Video, error)
	Remove(ctx context.Context, video *agg.Video) error
}
