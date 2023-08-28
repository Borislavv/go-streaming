package dtobuilder

import (
	"encoding/json"
	"github.com/Borislavv/video-streaming/internal/domain/dto"
	"net/http"
)

type VideoDtoBuilder struct {
}

func NewVideoDtoBuilder() *VideoDtoBuilder {
	return &VideoDtoBuilder{}
}

func (b *VideoDtoBuilder) BuildFromRequest(r *http.Request) (*dto.Video, error) {
	video := &dto.Video{}

	if err := json.NewDecoder(r.Body).Decode(&video); err != nil {
		return nil, err
	}

	return video, nil
}
