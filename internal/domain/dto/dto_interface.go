package dto

import (
	"github.com/Borislavv/video-streaming/internal/domain/vo"
	"mime/multipart"
	"time"
)

type CreateRequest interface {
	GetName() string
	GetResourceID() vo.ID
	GetDescription() string
}

type UpdateRequest interface {
	GetId() vo.ID
	GetName() string
	GetResourceID() vo.ID
	GetDescription() string
}

type DeleteRequest GetRequest
type GetRequest interface {
	GetId() vo.ID
}

type ListRequest interface {
	GetName() string         // path of name
	GetCreatedAt() time.Time // concrete search date point
	GetFrom() time.Time      // search date limit from
	GetTo() time.Time        // search date limit to
	PaginatedRequest
}

type UploadRequest interface {
	GetPart() *multipart.Part
	GetContentLength() int64
	GetUploadedFilename() string
	SetUploadedFilename(filename string)
	GetUploadedFilepath() string
	SetUploadedFilepath(filepath string)
	GetUploadedFilesize() int64
	SetUploadedFilesize(filesize int64)
	GetUploadedFiletype() string
	SetUploadedFiletype(filetype string)
}

type PaginatedRequest interface {
	GetPage() int
	GetLimit() int
}

type Resource interface {
	GetFilepath() string
	GetName() string
}

type Chunk interface {
	GetLen() int
	GetData() []byte
	SetData(data []byte)
	GetError() error
	SetError(err error)
}
