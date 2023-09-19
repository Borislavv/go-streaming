package agg

import (
	"github.com/Borislavv/video-streaming/internal/domain/entity"
	"github.com/Borislavv/video-streaming/internal/domain/vo"
)

type Audio struct {
	entity.Audio `bson:",inline"`
	Resource     entity.Resource `json:"resource" bson:"resource"`
	Timestamp    vo.Timestamp    `bson:",inline"`
}
