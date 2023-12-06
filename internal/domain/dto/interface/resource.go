package _interface

import (
	"github.com/Borislavv/video-streaming/internal/domain/vo"
	"net/http"
)

type Resource interface {
	GetFilepath() string
	GetName() string
}

type UploadResourceRequest interface {
	GetUserID() vo.ID
	GetRequest() *http.Request
	GetOriginFilename() string
	SetOriginFilename(filename string)
	GetUploadedFilename() string
	SetUploadedFilename(filename string)
	GetUploadedFilepath() string
	SetUploadedFilepath(filepath string)
	GetUploadedFilesize() int64
	SetUploadedFilesize(filesize int64)
	GetUploadedFiletype() string
	SetUploadedFiletype(filetype string)
}

type GetResourceRequest interface {
	GetID() vo.ID
	GetUserID() vo.ID
}

type DeleteResourceRequest GetResourceRequest
