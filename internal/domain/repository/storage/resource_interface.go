package repository

import (
	"context"
	"github.com/Borislavv/video-streaming/internal/domain/agg"
	"github.com/Borislavv/video-streaming/internal/infrastructure/repository/query"
)

type Resource interface {
	FindOneByID(ctx context.Context, q query.FindOneResourceByID) (*agg.Resource, error)
	Insert(ctx context.Context, resource *agg.Resource) (*agg.Resource, error)
	Remove(ctx context.Context, resource *agg.Resource) error
}
