package repository

import (
	"context"
	"github.com/Borislavv/video-streaming/internal/domain/agg"
	"github.com/Borislavv/video-streaming/internal/infrastructure/repository/query"
)

type Resource interface {
	FindOneByID(context.Context, query.FindOneResourceByID) (*agg.Resource, error)
	Insert(context.Context, *agg.Resource) (*agg.Resource, error)
	Remove(context.Context, *agg.Resource) error
}
