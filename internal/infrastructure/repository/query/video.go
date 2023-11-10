package query

import (
	"github.com/Borislavv/video-streaming/internal/domain/vo"
	"time"
)

type FindOneVideoByID interface {
	GetID() vo.ID
	GetUserID() vo.ID
}

type FindOneVideoByName interface {
	GetName() string
	GetUserID() vo.ID
}

type FindOneVideoByResourceID interface {
	GetResourceID() vo.ID
	GetUserID() vo.ID
}

type FindVideoList interface {
	GetName() string         // part of name
	GetUserID() vo.ID        // user identifier
	GetCreatedAt() time.Time // concrete search date point
	GetFrom() time.Time      // search date limit from
	GetTo() time.Time        // search date limit to
	Pagination
}
