package repository

import (
	"context"
	"github.com/Borislavv/video-streaming/internal/domain/agg"
	"github.com/Borislavv/video-streaming/internal/infrastructure/repository/query"
)

type BlockedToken interface {
	Insert(context.Context, *agg.BlockedToken) error
	Has(context.Context, query.HasBlockedToken) (found bool, err error)
}
