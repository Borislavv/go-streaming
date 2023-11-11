package entity

import (
	"github.com/Borislavv/video-streaming/internal/domain/vo"
	"time"
)

type BlockedToken struct {
	ID        vo.ID     `json:"id,omitempty" bson:"_id,omitempty,inline"`
	Value     string    `bson:"value"`
	UserID    vo.ID     `bson:"userID"`
	BlockedAt time.Time `bson:"blockedAt"`
}

func (r BlockedToken) GetID() vo.ID {
	return r.ID
}
