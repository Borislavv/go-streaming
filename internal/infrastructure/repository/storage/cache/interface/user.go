package cache_interface

import (
	"context"
	"github.com/Borislavv/video-streaming/internal/domain/agg"
	query_interface "github.com/Borislavv/video-streaming/internal/infrastructure/repository/query/interface"
)

type User interface {
	FindOneByID(context.Context, query_interface.FindOneUserByID) (*agg.User, error)
	FindOneByEmail(context.Context, query_interface.FindOneUserByEmail) (*agg.User, error)
	Insert(context.Context, *agg.User) (*agg.User, error)
	Update(context.Context, *agg.User) (*agg.User, error)
	Remove(context.Context, *agg.User) error
}
