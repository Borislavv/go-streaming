package user

import (
	"github.com/Borislavv/video-streaming/internal/domain/agg"
	"github.com/Borislavv/video-streaming/internal/domain/dto"
)

type CRUD interface {
	Get(reqDTO dto.GetUserRequest) (*agg.User, error)
	Create(reqDTO dto.CreateUserRequest) (*agg.User, error)
	//Update(reqDTO dto.UpdateVideoRequest) (*agg.User, error)
	//Delete(reqDTO dto.DeleteVideoRequest) error
}
