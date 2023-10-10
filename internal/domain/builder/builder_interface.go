package builder

import (
	"github.com/Borislavv/video-streaming/internal/domain/agg"
	"github.com/Borislavv/video-streaming/internal/domain/dto"
	"net/http"
)

type Resource interface {
	BuildUploadRequestDTOFromRequest(r *http.Request) (*dto.ResourceUploadRequestDTO, error)
	BuildAggFromUploadRequestDTO(req dto.UploadRequest) *agg.Resource
}

type Video interface {
	BuildGetRequestDTOFromRequest(r *http.Request) (*dto.VideoGetRequestDTO, error)
	BuildListRequestDTOFromRequest(r *http.Request) (*dto.VideoListRequestDTO, error)
	BuildCreateRequestDTOFromRequest(r *http.Request) (*dto.VideoCreateRequestDTO, error)
	BuildAggFromCreateRequestDTO(dto dto.CreateRequest) (*agg.Video, error)
	BuildUpdateRequestDTOFromRequest(r *http.Request) (*dto.VideoUpdateRequestDTO, error)
	BuildAggFromUpdateRequestDTO(dto dto.UpdateRequest) (*agg.Video, error)
	BuildDeleteRequestDTOFromRequest(r *http.Request) (*dto.VideoDeleteRequestDto, error)
}

type User interface {
	BuildGetRequestDTOFromRequest(r *http.Request) (*dto.UserGetRequestDTO, error)
}
