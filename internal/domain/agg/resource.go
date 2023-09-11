package agg

import (
	"github.com/Borislavv/video-streaming/internal/domain/entity"
	"github.com/Borislavv/video-streaming/internal/domain/vo"
)

type Resource struct {
	ID        vo.ID           `json:"id" bson:",inline"`
	Resource  entity.Resource `json:"resource" bson:",inline"`
	Timestamp vo.Timestamp    `json:"timestamp" bson:",inline"`
}
