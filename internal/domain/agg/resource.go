package agg

import (
	"github.com/Borislavv/video-streaming/internal/domain/entity"
	"github.com/Borislavv/video-streaming/internal/domain/vo"
)

type Resource struct {
	entity.Resource `bson:",inline"`
	Timestamp       vo.Timestamp `json:"timestamp" bson:",inline"`
}

func (r *Resource) GetName() string {
	return r.Name
}
func (r *Resource) GetFilepath() string {
	return r.Filepath
}
