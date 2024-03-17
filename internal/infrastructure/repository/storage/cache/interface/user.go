package cacheinterface

import (
	"context"
	"github.com/Borislavv/video-streaming/internal/domain/agg"
	queryinterface "github.com/Borislavv/video-streaming/internal/infrastructure/repository/query/interface"
)

type User interface {
	FindOneByID(context.Context, queryinterface.FindOneUserByID) (*agg.User, error)
	FindOneByEmail(context.Context, queryinterface.FindOneUserByEmail) (*agg.User, error)
	Insert(context.Context, *agg.User) (*agg.User, error)
	Update(context.Context, *agg.User) (*agg.User, error)
	Remove(context.Context, *agg.User) error
}
