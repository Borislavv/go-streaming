package agg

import (
	"github.com/Borislavv/video-streaming/internal/domain/entity"
	"github.com/Borislavv/video-streaming/internal/domain/vo"
)

type Video struct {
	ID        vo.ID        `json:"id" bson:",inline"`
	Video     entity.Video `json:"video" bson:",inline"`
	Timestamp vo.Timestamp `json:"timestamp" bson:",inline"`
}
