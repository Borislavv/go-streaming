package _interface

import (
	"github.com/Borislavv/video-streaming/internal/domain/agg"
	"github.com/Borislavv/video-streaming/internal/domain/dto"
)

type CRUD interface {
	Get(reqDTO dto.GetUserRequest) (*agg.User, error)
	Create(reqDTO dto.CreateUserRequest) (*agg.User, error)
	Update(reqDTO dto.UpdateUserRequest) (*agg.User, error)
	Delete(reqDTO dto.DeleteUserRequest) error
}
