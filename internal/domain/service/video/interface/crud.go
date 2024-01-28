package videointerface

import (
	"github.com/Borislavv/video-streaming/internal/domain/agg"
	dtointerface "github.com/Borislavv/video-streaming/internal/domain/dto/interface"
)

type CRUD interface {
	Get(reqDTO dtointerface.GetVideoRequest) (*agg.Video, error)
	List(reqDTO dtointerface.ListVideoRequest) (list []*agg.Video, total int64, err error)
	Create(reqDTO dtointerface.CreateVideoRequest) (*agg.Video, error)
	Update(reqDTO dtointerface.UpdateVideoRequest) (*agg.Video, error)
	Delete(reqDTO dtointerface.DeleteVideoRequest) error
}
