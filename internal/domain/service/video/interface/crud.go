package video_interface

import (
	"github.com/Borislavv/video-streaming/internal/domain/agg"
	"github.com/Borislavv/video-streaming/internal/domain/dto/interface"
)

type CRUD interface {
	Get(reqDTO dto_interface.GetVideoRequest) (*agg.Video, error)
	List(reqDTO dto_interface.ListVideoRequest) (list []*agg.Video, total int64, err error)
	Create(reqDTO dto_interface.CreateVideoRequest) (*agg.Video, error)
	Update(reqDTO dto_interface.UpdateVideoRequest) (*agg.Video, error)
	Delete(reqDTO dto_interface.DeleteVideoRequest) error
}
