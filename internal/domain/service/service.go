package service

import (
	"github.com/Borislavv/video-streaming/internal/domain/agg"
	"github.com/Borislavv/video-streaming/internal/domain/dto"
)

type Resource interface {
	Upload(req dto.UploadRequest) (*agg.Resource, error)
}

type Video interface {
	Get(req dto.GetRequest) (*agg.Video, error)
	List(req dto.ListRequest) (list []*agg.Video, total int64, err error)
	Create(req dto.CreateRequest) (*agg.Video, error)
	Update(req dto.UpdateRequest) (*agg.Video, error)
	Delete(req dto.DeleteRequest) error
}
