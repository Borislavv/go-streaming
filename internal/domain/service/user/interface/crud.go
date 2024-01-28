package userinterface

import (
	"github.com/Borislavv/video-streaming/internal/domain/agg"
	"github.com/Borislavv/video-streaming/internal/domain/dto/interface"
)

type CRUD interface {
	Get(reqDTO dtointerface.GetUserRequest) (*agg.User, error)
	Create(reqDTO dtointerface.CreateUserRequest) (*agg.User, error)
	Update(reqDTO dtointerface.UpdateUserRequest) (*agg.User, error)
	Delete(reqDTO dtointerface.DeleteUserRequest) error
}
