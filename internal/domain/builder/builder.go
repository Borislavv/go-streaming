package builder

import (
	"github.com/Borislavv/video-streaming/internal/domain/agg"
	"github.com/Borislavv/video-streaming/internal/domain/dto"
	"net/http"
)

type Video interface {
	BuildGetRequestDtoFromRequest(r *http.Request) (*dto.VideoGetRequestDto, error)
	BuildListRequestDtoFromRequest(r *http.Request) (*dto.VideoListRequestDto, error)
	BuildCreateRequestDtoFromRequest(r *http.Request) (*dto.VideoCreateRequestDto, error)
	BuildUpdateRequestDtoFromRequest(r *http.Request) (*dto.VideoUpdateRequestDto, error)
	BuildDeleteRequestDtoFromRequest(r *http.Request) (*dto.VideoDeleteRequestDto, error)

	BuildAggFromCreateRequestDto(dto dto.CreateRequest) *agg.Video
	BuildAggFromUpdateRequestDto(dto dto.UpdateRequest) (*agg.Video, error)
}
