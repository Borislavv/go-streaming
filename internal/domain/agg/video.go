package agg

import (
	"github.com/Borislavv/video-streaming/internal/domain/entity"
	"github.com/Borislavv/video-streaming/internal/domain/vo"
)

type Video struct {
	ID        vo.ID        `bson:",inline"`
	Video     entity.Video `bson:",inline"`
	Timestamp vo.Timestamp `bson:",inline"`
}
