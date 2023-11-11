package repository

import (
	"context"
	"github.com/Borislavv/video-streaming/internal/domain/agg"
	"github.com/Borislavv/video-streaming/internal/domain/dto"
	"github.com/Borislavv/video-streaming/internal/domain/vo"
)

type Audio interface {
	Insert(ctx context.Context, audio *agg.Audio) (string, error)
	Update(ctx context.Context, audio *agg.Audio) error
	Find(ctx context.Context, id vo.ID) (*agg.Audio, error)
	FindList(ctx context.Context, dto dto.ListAudioRequest) ([]*agg.Audio, error)
}
