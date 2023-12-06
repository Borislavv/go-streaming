package _interface

import (
	"context"
	"github.com/Borislavv/video-streaming/internal/domain/agg"
	"github.com/Borislavv/video-streaming/internal/infrastructure/repository/query"
)

type User interface {
	FindOneByID(context.Context, query.FindOneUserByID) (*agg.User, error)
	FindOneByEmail(context.Context, query.FindOneUserByEmail) (*agg.User, error)
	Insert(context.Context, *agg.User) (*agg.User, error)
	Update(context.Context, *agg.User) (*agg.User, error)
	Remove(context.Context, *agg.User) error
}
