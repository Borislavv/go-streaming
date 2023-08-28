package entitybuilder

import (
	"github.com/Borislavv/video-streaming/internal/domain/dto"
	"github.com/Borislavv/video-streaming/internal/domain/entity"
)

type VideoEntityBuilder struct {
}

func NewVideoEntityBuilder() *VideoEntityBuilder {
	return &VideoEntityBuilder{}
}

func (b *VideoEntityBuilder) Build(video *dto.Video) *entity.Video {
	return &entity.Video{
		Name:        video.Name,
		Path:        video.Path,
		Description: video.Description,
	}
}
