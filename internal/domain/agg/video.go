package agg

import (
	"github.com/Borislavv/video-streaming/internal/domain/entity"
	"github.com/Borislavv/video-streaming/internal/domain/vo"
)

type Video struct {
	entity.Video `bson:",inline"`
	Resource     entity.Resource `json:"resource" bson:"resource"`
	Timestamp    vo.Timestamp    `json:"timestamp" bson:",inline"`
}
