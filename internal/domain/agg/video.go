package agg

import (
	"github.com/Borislavv/video-streaming/internal/domain/entity"
	"github.com/Borislavv/video-streaming/internal/domain/vo"
)

type Video struct {
	ID           vo.ID `json:"id" bson:",inline"`
	entity.Video `json:"video" bson:",inline"`
	Resource     entity.Resource `json:"resource" bson:"resource"`
	Timestamp    vo.Timestamp    `json:"timestamp" bson:",inline"`
}
