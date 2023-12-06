package _interface

import (
	"github.com/Borislavv/video-streaming/internal/domain/vo"
	"time"
)

type CreateAudioRequest interface {
	GetName() string
	GetResourceID() vo.ID
	GetDescription() string
}

type UpdateAudioRequest interface {
	GetID() vo.ID
	GetName() string
	GetResourceID() vo.ID
	GetDescription() string
}

type GetAudioRequest interface {
	GetID() vo.ID
}
type ListAudioRequest interface {
	GetName() string         // path of name
	GetCreatedAt() time.Time // concrete search date point
	GetFrom() time.Time      // search date limit from
	GetTo() time.Time        // search date limit to
	PaginatedRequest
}

type DeleteAudioRequest GetAudioRequest
