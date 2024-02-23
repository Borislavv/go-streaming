package dtointerface

import (
	"github.com/Borislavv/video-streaming/internal/domain/vo"
	"time"
)

type CreateVideoRequest interface {
	GetName() string
	GetUserID() vo.ID
	GetResourceID() vo.ID
	GetDescription() string
}

type UpdateVideoRequest interface {
	GetID() vo.ID
	GetName() string
	GetUserID() vo.ID
	GetResourceID() vo.ID
	GetDescription() string
}

type GetVideoRequest interface {
	GetID() vo.ID
	GetUserID() vo.ID
}
type ListVideoRequest interface {
	GetName() string         // part of name
	GetUserID() vo.ID        // user identifier
	GetCreatedAt() time.Time // concrete search date point
	GetFrom() time.Time      // search date limit from
	GetTo() time.Time        // search date limit to
	PaginatedRequest
}

type DeleteVideoRequest GetVideoRequest
