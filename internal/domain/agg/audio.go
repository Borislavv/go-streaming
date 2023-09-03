package agg

import (
	"github.com/Borislavv/video-streaming/internal/domain/entity"
	"github.com/Borislavv/video-streaming/internal/domain/vo"
)

type Audio struct {
	ID        vo.ID        `json:"id,omitempty" bson:"_id,omitempty,inline"`
	Audio     entity.Audio `bson:",inline"`
	Timestamp vo.Timestamp `bson:",inline"`
}
