package video

import (
	"github.com/Borislavv/video-streaming/internal/domain/agg"
	"github.com/Borislavv/video-streaming/internal/domain/dto"
)

type CRUD interface {
	Get(reqDTO dto.GetVideoRequest) (*agg.Video, error)
	List(reqDTO dto.ListVideoRequest) (list []*agg.Video, total int64, err error)
	Create(reqDTO dto.CreateVideoRequest) (*agg.Video, error)
	Update(reqDTO dto.UpdateVideoRequest) (*agg.Video, error)
	Delete(reqDTO dto.DeleteVideoRequest) error
}
