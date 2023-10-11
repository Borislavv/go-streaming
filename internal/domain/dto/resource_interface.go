package dto

import "net/http"

type Resource interface {
	GetFilepath() string
	GetName() string
}

type UploadResourceRequest interface {
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
