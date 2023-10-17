package builder

import (
	"github.com/Borislavv/video-streaming/internal/domain/agg"
	"github.com/Borislavv/video-streaming/internal/domain/dto"
	"net/http"
)

type User interface {
	BuildGetRequestDTOFromRequest(r *http.Request) (*dto.UserGetRequestDTO, error)
	BuildCreateRequestDTOFromRequest(r *http.Request) (*dto.UserCreateRequestDTO, error)
	BuildAggFromCreateRequestDTO(dto dto.CreateUserRequest) (*agg.User, error)
	BuildUpdateRequestDTOFromRequest(r *http.Request) (*dto.UserUpdateRequestDTO, error)
	BuildAggFromUpdateRequestDTO(dto dto.UpdateUserRequest) (*agg.User, error)
	BuildDeleteRequestDTOFromRequest(r *http.Request) (*dto.UserDeleteRequestDto, error)
}
