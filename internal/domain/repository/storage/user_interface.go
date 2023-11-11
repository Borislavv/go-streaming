package repository

import (
	"context"
	"github.com/Borislavv/video-streaming/internal/domain/agg"
	"github.com/Borislavv/video-streaming/internal/domain/vo"
)

type User interface {
	FindOneByID(ctx context.Context, id vo.ID) (*agg.User, error)
	FindOneByEmail(ctx context.Context, email string) (*agg.User, error)
	Insert(ctx context.Context, user *agg.User) (*agg.User, error)
	Update(ctx context.Context, user *agg.User) (*agg.User, error)
	Remove(ctx context.Context, user *agg.User) error
}
