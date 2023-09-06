package repository

import (
	"context"
	"github.com/Borislavv/video-streaming/internal/domain/agg"
	"github.com/Borislavv/video-streaming/internal/domain/vo"
)

type Video interface {
	Find(ctx context.Context, id vo.ID) (*agg.Video, error)
	Insert(ctx context.Context, video *agg.Video) (*agg.Video, error)
	Update(ctx context.Context, video *agg.Video) (*agg.Video, error)
	//FindBy(ctx context.Context, id vo.ID) (*agg.Video, error)
	//FindList(ctx context.Context, dto dto.ListRequest) ([]*agg.Video, error)
}
