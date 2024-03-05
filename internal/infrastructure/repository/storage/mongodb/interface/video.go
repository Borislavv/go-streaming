package mongodbinterface

import (
	"context"
	"github.com/Borislavv/video-streaming/internal/domain/agg"
	queryinterface "github.com/Borislavv/video-streaming/internal/infrastructure/repository/query/interface"
)

type Video interface {
	FindOneByID(ctx context.Context, q queryinterface.FindOneVideoByID) (*agg.Video, error)
	FindOneByName(ctx context.Context, q queryinterface.FindOneVideoByName) (*agg.Video, error)
	FindOneByResourceID(ctx context.Context, q queryinterface.FindOneVideoByResourceID) (*agg.Video, error)
	FindList(ctx context.Context, q queryinterface.FindVideoList) (list []*agg.Video, total int64, err error)
	Insert(ctx context.Context, video *agg.Video) (*agg.Video, error)
	Update(ctx context.Context, video *agg.Video) (*agg.Video, error)
	Remove(ctx context.Context, video *agg.Video) error
}
