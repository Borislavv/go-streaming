package video

import (
	"github.com/Borislavv/video-streaming/internal/domain/agg"
	"github.com/Borislavv/video-streaming/internal/domain/dto"
)

type CRUD interface {
	Get(reqDTO dto.GetRequest) (*agg.Video, error)
	List(reqDTO dto.ListRequest) (list []*agg.Video, total int64, err error)
	Create(reqDTO dto.CreateRequest) (*agg.Video, error)
	Update(reqDTO dto.UpdateRequest) (*agg.Video, error)
	Delete(reqDTO dto.DeleteRequest) error
}
