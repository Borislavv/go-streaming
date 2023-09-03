package builder

import (
	"github.com/Borislavv/video-streaming/internal/domain/agg"
	"github.com/Borislavv/video-streaming/internal/domain/dto"
	"net/http"
)

type Video interface {
	BuildCreateRequestDtoFromRequest(r *http.Request) (*dto.VideoCreateRequestDto, error)
	BuildAggFromCreateRequestDto(dto dto.CreateRequest) *agg.Video

	BuildUpdateRequestDtoFromRequest(r *http.Request) (*dto.VideoUpdateRequestDto, error)
	BuildAggFromUpdateRequestDto(dto dto.UpdateRequest) (*agg.Video, error)

	BuildListRequestDtoFromRequest(r *http.Request) (*dto.VideoListRequestDto, error)
}
