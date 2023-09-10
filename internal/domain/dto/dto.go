package dto

import "github.com/Borislavv/video-streaming/internal/domain/vo"

type CreateRequest interface {
	GetName() string
	GetPath() string
	GetDescription() string
}

type UpdateRequest interface {
	GetId() vo.ID
	GetName() string
	GetDescription() string
}

type DeleteRequest GetRequest
type GetRequest interface {
	GetId() vo.ID
}

type ListRequest interface {
	GetName() string // path of name
	GetPath() string // path of path
	PaginatedRequest
}

type PaginatedRequest interface {
	GetPage() int
	GetLimit() int
}

type Chunk interface {
	GetLen() int
	SetLen(len int)
	GetData() []byte
	SetData(data []byte)
	GetError() error
	SetError(err error)
}
