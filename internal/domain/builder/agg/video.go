package aggbuilder

import (
	"github.com/Borislavv/video-streaming/data/vo"
	"github.com/Borislavv/video-streaming/internal/domain/agg"
	"github.com/Borislavv/video-streaming/internal/domain/entity"
	"time"
)

type VideoAggBuilder struct {
}

func NewVideoAggBuilder() *VideoAggBuilder {
	return &VideoAggBuilder{}
}

func (b *VideoAggBuilder) Build(video *entity.Video) *agg.Video {
	return &agg.Video{
		Video: *video,
		Timestamp: vo.Timestamp{
			CreatedAt: time.Now(),
		},
	}
}
