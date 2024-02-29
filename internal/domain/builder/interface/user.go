package builderinterface

import (
	"github.com/Borislavv/video-streaming/internal/domain/agg"
	"github.com/Borislavv/video-streaming/internal/domain/dto"
	dtointerface "github.com/Borislavv/video-streaming/internal/domain/dto/interface"
	"net/http"
)

type User interface {
	BuildGetRequestDTOFromRequest(r *http.Request) (*dto.UserGetRequestDTO, error)
	BuildCreateRequestDTOFromRequest(r *http.Request) (*dto.UserCreateRequestDTO, error)
	BuildAggFromCreateRequestDTO(reqDTO dtointerface.CreateUserRequest) (*agg.User, error)
	BuildUpdateRequestDTOFromRequest(r *http.Request) (*dto.UserUpdateRequestDTO, error)
	BuildAggFromUpdateRequestDTO(reqDTO dtointerface.UpdateUserRequest) (*agg.User, error)
	BuildDeleteRequestDTOFromRequest(r *http.Request) (*dto.UserDeleteRequestDTO, error)
	BuildResponseDTO(user *agg.User) (*dto.UserResponseDTO, error)
}
