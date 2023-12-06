package _interface

import (
	"context"
	"github.com/Borislavv/video-streaming/internal/domain/agg"
)

type BlockedToken interface {
	Insert(ctx context.Context, token *agg.BlockedToken) error
	Has(ctx context.Context, token string) (found bool, err error)
}
