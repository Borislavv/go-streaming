package user

import (
	"github.com/Borislavv/video-streaming/internal/domain/agg"
	"github.com/Borislavv/video-streaming/internal/domain/dto"
)

type CRUD interface {
	Get(reqDTO dto.GetRequest) (*agg.User, error)
	//Create(reqDTO dto.CreateRequest) (*agg.User, error)
	//Update(reqDTO dto.UpdateRequest) (*agg.User, error)
	//Delete(reqDTO dto.DeleteRequest) error
}
