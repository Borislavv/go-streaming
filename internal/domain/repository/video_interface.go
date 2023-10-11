package repository

import (
	"context"
	"github.com/Borislavv/video-streaming/internal/domain/agg"
	"github.com/Borislavv/video-streaming/internal/domain/dto"
	"github.com/Borislavv/video-streaming/internal/domain/vo"
)

type Video interface {
	Find(ctx context.Context, id vo.ID) (*agg.Video, error)
	FindList(ctx context.Context, dto dto.ListVideoRequest) (list []*agg.Video, total int64, err error)
	FindOneByName(ctx context.Context, name string) (*agg.Video, error)
	FindOneByResourceId(ctx context.Context, resourceID vo.ID) (*agg.Video, error)
	Insert(ctx context.Context, video *agg.Video) (*agg.Video, error)
	Update(ctx context.Context, video *agg.Video) (*agg.Video, error)
	Remove(ctx context.Context, video *agg.Video) error
}
