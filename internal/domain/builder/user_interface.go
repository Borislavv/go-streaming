package builder

import (
	"github.com/Borislavv/video-streaming/internal/domain/agg"
	"github.com/Borislavv/video-streaming/internal/domain/dto"
	"net/http"
)

type User interface {
	BuildGetRequestDTOFromRequest(r *http.Request) (*dto.UserGetRequestDTO, error)
	BuildCreateRequestDTOFromRequest(r *http.Request) (*dto.UserCreateRequestDTO, error)
	BuildAggFromCreateRequestDTO(reqDTO dto.CreateUserRequest) (*agg.User, error)
	BuildUpdateRequestDTOFromRequest(r *http.Request) (*dto.UserUpdateRequestDTO, error)
	BuildAggFromUpdateRequestDTO(reqDTO dto.UpdateUserRequest) (*agg.User, error)
	BuildDeleteRequestDTOFromRequest(r *http.Request) (*dto.UserDeleteRequestDTO, error)
	BuildResponseDTO(user *agg.User) (*dto.UserResponseDTO, error)
}
