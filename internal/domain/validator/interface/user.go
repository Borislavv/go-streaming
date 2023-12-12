package validator_interface

import (
	"github.com/Borislavv/video-streaming/internal/domain/agg"
	dto_interface "github.com/Borislavv/video-streaming/internal/domain/dto/interface"
)

type User interface {
	ValidateGetRequestDTO(reqDTO dto_interface.GetUserRequest) error
	ValidateCreateRequestDTO(reqDTO dto_interface.CreateUserRequest) error
	ValidateUpdateRequestDTO(reqDTO dto_interface.UpdateUserRequest) error
	ValidateDeleteRequestDTO(reqDTO dto_interface.DeleteUserRequest) error
	ValidateAggregate(agg *agg.User) error
}
