package aggbuilder

import (
	"github.com/Borislavv/video-streaming/internal/domain/agg"
	"github.com/Borislavv/video-streaming/internal/domain/dto"
)

type Video interface {
	Build(video *dto.Video) *agg.Video
}
