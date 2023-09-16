package agg

import (
	"github.com/Borislavv/video-streaming/internal/domain/entity"
	"github.com/Borislavv/video-streaming/internal/domain/vo"
)

type Resource struct {
	entity.Resource `json:"resource" bson:",inline"`
	Timestamp       vo.Timestamp `json:"timestamp" bson:",inline"`
}

func (r *Resource) GetFilepath() string {
	return r.Filepath
}
