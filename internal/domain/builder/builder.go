package builder

import (
	"github.com/Borislavv/video-streaming/internal/domain/agg"
	"github.com/Borislavv/video-streaming/internal/domain/dto"
	"net/http"
)

type Resource interface {
	BuildUploadRequestDtoFromRequest(r *http.Request) (*dto.ResourceUploadRequestDTO, error)
	BuildAggFromUploadRequestDto(req dto.UploadRequest) *agg.Resource
}

type Video interface {
	BuildGetRequestDtoFromRequest(r *http.Request) (*dto.VideoGetRequestDTO, error)
	BuildListRequestDtoFromRequest(r *http.Request) (*dto.VideoListRequestDTO, error)
	BuildCreateRequestDtoFromRequest(r *http.Request) (*dto.VideoCreateRequestDTO, error)
	BuildAggFromCreateRequestDto(dto dto.CreateRequest) (*agg.Video, error)
	BuildUpdateRequestDtoFromRequest(r *http.Request) (*dto.VideoUpdateRequestDTO, error)
	BuildAggFromUpdateRequestDto(dto dto.UpdateRequest) (*agg.Video, error)
	BuildDeleteRequestDtoFromRequest(r *http.Request) (*dto.VideoDeleteRequestDto, error)
}
