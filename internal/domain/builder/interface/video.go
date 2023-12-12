package builder_interface

import (
	"github.com/Borislavv/video-streaming/internal/domain/agg"
	"github.com/Borislavv/video-streaming/internal/domain/dto"
	dto_interface "github.com/Borislavv/video-streaming/internal/domain/dto/interface"
	"net/http"
)

type Video interface {
	BuildGetRequestDTOFromRequest(r *http.Request) (*dto.VideoGetRequestDTO, error)
	BuildListRequestDTOFromRequest(r *http.Request) (*dto.VideoListRequestDTO, error)
	BuildCreateRequestDTOFromRequest(r *http.Request) (*dto.VideoCreateRequestDTO, error)
	BuildAggFromCreateRequestDTO(reqDTO dto_interface.CreateVideoRequest) (*agg.Video, error)
	BuildUpdateRequestDTOFromRequest(r *http.Request) (*dto.VideoUpdateRequestDTO, error)
	BuildAggFromUpdateRequestDTO(reqDTO dto_interface.UpdateVideoRequest) (*agg.Video, error)
	BuildDeleteRequestDTOFromRequest(r *http.Request) (*dto.VideoDeleteRequestDto, error)
}
