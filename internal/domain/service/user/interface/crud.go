package user_interface

import (
	"github.com/Borislavv/video-streaming/internal/domain/agg"
	"github.com/Borislavv/video-streaming/internal/domain/dto/interface"
)

type CRUD interface {
	Get(reqDTO dto_interface.GetUserRequest) (*agg.User, error)
	Create(reqDTO dto_interface.CreateUserRequest) (*agg.User, error)
	Update(reqDTO dto_interface.UpdateUserRequest) (*agg.User, error)
	Delete(reqDTO dto_interface.DeleteUserRequest) error
}
