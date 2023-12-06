package _interface

import (
	"github.com/Borislavv/video-streaming/internal/domain/agg"
	"github.com/Borislavv/video-streaming/internal/domain/dto"
)

type User interface {
	ValidateGetRequestDTO(reqDTO dto.GetUserRequest) error
	ValidateCreateRequestDTO(reqDTO dto.CreateUserRequest) error
	ValidateUpdateRequestDTO(reqDTO dto.UpdateUserRequest) error
	ValidateDeleteRequestDTO(reqDTO dto.DeleteUserRequest) error
	ValidateAggregate(agg *agg.User) error
}
