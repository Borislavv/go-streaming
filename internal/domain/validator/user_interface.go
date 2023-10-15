package validator

import (
	"github.com/Borislavv/video-streaming/internal/domain/agg"
	"github.com/Borislavv/video-streaming/internal/domain/dto"
)

type User interface {
	ValidateGetRequestDTO(req dto.GetUserRequest) error
	ValidateCreateRequestDTO(req dto.CreateUserRequest) error
	//ValidateUpdateRequestDTO(req dto.UpdateUserRequest) error
	ValidateDeleteRequestDTO(req dto.DeleteUserRequest) error
	ValidateAggregate(agg *agg.User) error
}
