package repository

import (
	"context"
	"github.com/Borislavv/video-streaming/internal/domain/agg"
)

type Video interface {
	Insert(ctx context.Context, video *agg.Video) (string, error)
	//Update(ctx context.Context, video *agg.Video) error
	//Find(ctx context.Context, id vo.ID) (*agg.Video, error)
	//FindList(ctx context.Context, dto dto.ListRequest) ([]*agg.Video, error)
}
