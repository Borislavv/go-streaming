package aggbuilder

import (
	"github.com/Borislavv/video-streaming/internal/domain/agg"
	"github.com/Borislavv/video-streaming/internal/domain/entity"
)

type Video interface {
	Build(video *entity.Video) *agg.Video
}
