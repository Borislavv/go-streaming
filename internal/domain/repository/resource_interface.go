package repository

import (
	"context"
	"github.com/Borislavv/video-streaming/internal/domain/agg"
	"github.com/Borislavv/video-streaming/internal/domain/vo"
)

type Resource interface {
	Find(ctx context.Context, id vo.ID) (*agg.Resource, error)
	Insert(ctx context.Context, resource *agg.Resource) (*agg.Resource, error)
	Remove(ctx context.Context, resource *agg.Resource) error
}
