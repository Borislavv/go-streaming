package validatorinterface

import (
	"github.com/Borislavv/video-streaming/internal/domain/agg"
	dtointerface "github.com/Borislavv/video-streaming/internal/domain/dto/interface"
)

type User interface {
	ValidateGetRequestDTO(reqDTO dtointerface.GetUserRequest) error
	ValidateCreateRequestDTO(reqDTO dtointerface.CreateUserRequest) error
	ValidateUpdateRequestDTO(reqDTO dtointerface.UpdateUserRequest) error
	ValidateDeleteRequestDTO(reqDTO dtointerface.DeleteUserRequest) error
	ValidateAggregate(agg *agg.User) error
}
