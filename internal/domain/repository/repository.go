package repository

import (
	"context"
	"github.com/Borislavv/video-streaming/internal/domain/agg"
)

type VideoRepository interface {
	Insert(ctx context.Context, video *agg.Video) (string, error)
}
