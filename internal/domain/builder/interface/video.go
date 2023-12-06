package _interface

import (
	"github.com/Borislavv/video-streaming/internal/domain/agg"
	"github.com/Borislavv/video-streaming/internal/domain/dto"
	"net/http"
)

type Video interface {
	BuildGetRequestDTOFromRequest(r *http.Request) (*dto.VideoGetRequestDTO, error)
	BuildListRequestDTOFromRequest(r *http.Request) (*dto.VideoListRequestDTO, error)
	BuildCreateRequestDTOFromRequest(r *http.Request) (*dto.VideoCreateRequestDTO, error)
	BuildAggFromCreateRequestDTO(reqDTO dto.CreateVideoRequest) (*agg.Video, error)
	BuildUpdateRequestDTOFromRequest(r *http.Request) (*dto.VideoUpdateRequestDTO, error)
	BuildAggFromUpdateRequestDTO(reqDTO dto.UpdateVideoRequest) (*agg.Video, error)
	BuildDeleteRequestDTOFromRequest(r *http.Request) (*dto.VideoDeleteRequestDto, error)
}
