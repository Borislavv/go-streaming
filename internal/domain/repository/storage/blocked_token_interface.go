package repository

import (
	"context"
	"github.com/Borislavv/video-streaming/internal/domain/agg"
)

type BlockedToken interface {
	Insert(context.Context, *agg.BlockedToken) error
	Has(ctx context.Context, token string) (found bool, err error)
}
