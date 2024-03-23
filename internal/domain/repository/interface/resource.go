package repositoryinterface

import (
	"context"
	"github.com/Borislavv/video-streaming/internal/domain/agg"
)

type Resource interface {
	FindOneByID(context.Context, queryinterface.FindOneResourceByID) (*agg.Resource, error)
	Insert(context.Context, *agg.Resource) (*agg.Resource, error)
	Remove(context.Context, *agg.Resource) error
}
