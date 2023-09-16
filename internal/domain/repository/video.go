package repository

import (
	"context"
	"github.com/Borislavv/video-streaming/internal/domain/agg"
	"github.com/Borislavv/video-streaming/internal/domain/dto"
	"github.com/Borislavv/video-streaming/internal/domain/vo"
)

type Video interface {
	Find(ctx context.Context, id vo.ID) (*agg.Video, error)
	FindByResource(ctx context.Context, resource *agg.Resource) (*agg.Video, error)
	FindList(ctx context.Context, query dto.ListRequest) ([]*agg.Video, error)
	Insert(ctx context.Context, video *agg.Video) (*agg.Video, error)
	Update(ctx context.Context, video *agg.Video) (*agg.Video, error)
	Has(ctx context.Context, video *agg.Video) (bool, error)
	Remove(ctx context.Context, video *agg.Video) error
}
