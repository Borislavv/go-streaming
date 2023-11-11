package repository

import (
	"context"
)

type BlockedToken interface {
	Insert(ctx context.Context, token string) error
	Has(ctx context.Context, token string) (found bool, err error)
}
