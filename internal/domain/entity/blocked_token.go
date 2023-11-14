package entity

import (
	"github.com/Borislavv/video-streaming/internal/domain/vo"
	"time"
)

type BlockedToken struct {
	ID        vo.ID     `json:"id,omitempty" bson:"_id,omitempty,inline"`
	UserID    vo.ID     `bson:"user"`
	Value     string    `bson:"value"`
	Reason    string    `bson:"reason"`
	BlockedAt time.Time `bson:"blockedAt"`
}

func (r BlockedToken) GetID() vo.ID {
	return r.ID
}
