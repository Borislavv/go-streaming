package agg

import (
	"github.com/Borislavv/video-streaming/internal/domain/entity"
	"github.com/Borislavv/video-streaming/internal/domain/vo"
)

type User struct {
	entity.User `bson:",inline"`

	Timestamp vo.Timestamp `json:"timestamp" bson:",inline"`
}
