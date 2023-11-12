package agg

import (
	"github.com/Borislavv/video-streaming/internal/domain/entity"
	"github.com/Borislavv/video-streaming/internal/domain/vo"
	"time"
)

type BlockedToken struct {
	entity.BlockedToken `bson:",inline"`

	Timestamp vo.Timestamp `json:"timestamp" bson:",inline"`
}

func NewBlockedToken(token string, userID vo.ID) *BlockedToken {
	return &BlockedToken{
		BlockedToken: entity.BlockedToken{
			Value:     token,
			UserID:    userID,
			BlockedAt: time.Now(),
		},
		Timestamp: vo.Timestamp{
			CreatedAt: time.Now(),
			UpdatedAt: time.Time{},
		},
	}
}
