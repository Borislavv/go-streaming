package agg

import (
	"github.com/Borislavv/video-streaming/internal/domain/entity"
	"github.com/Borislavv/video-streaming/internal/domain/vo"
)

type BlockedToken struct {
	entity.BlockedToken `bson:",inline"`

	Timestamp vo.Timestamp `json:"timestamp" bson:",inline"`
}
