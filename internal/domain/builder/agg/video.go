package aggbuilder

import (
	"github.com/Borislavv/video-streaming/data/vo"
	"github.com/Borislavv/video-streaming/internal/domain/agg"
	"github.com/Borislavv/video-streaming/internal/domain/dto"
	"github.com/Borislavv/video-streaming/internal/domain/entity"
	"time"
)

type VideoAggBuilder struct {
}

func NewVideoAggBuilder() *VideoAggBuilder {
	return &VideoAggBuilder{}
}

func (b *VideoAggBuilder) Build(video *dto.Video) *agg.Video {
	return &agg.Video{
		Video: entity.Video{
			Name:        video.Name,
			Path:        video.Path,
			Description: video.Description,
		},
		Timestamp: vo.Timestamp{
			CreatedAt: time.Now(),
		},
	}
}
