package entitybuilder

import (
	"github.com/Borislavv/video-streaming/internal/domain/dto"
	"github.com/Borislavv/video-streaming/internal/domain/entity"
)

type Video interface {
	Build(video *dto.Video) *entity.Video
}
