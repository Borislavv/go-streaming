package dto

import (
	"net/http"
)

type ResourceUploadRequestDTO struct {
	request          *http.Request
	contentLength    int64
	uploadedFilename string
	uploadedFilepath string
	uploadedFiletype string
	uploadedFilesize int64
}

func NewResourceUploadRequest(r *http.Request) (dto *ResourceUploadRequestDTO) {
	return &ResourceUploadRequestDTO{request: r}
}

func (r *ResourceUploadRequestDTO) GetRequest() *http.Request {
	return r.request
}
func (r *ResourceUploadRequestDTO) GetUploadedFilename() string {
	return r.uploadedFilename
}
func (r *ResourceUploadRequestDTO) SetUploadedFilename(filename string) {
	r.uploadedFilename = filename
}
func (r *ResourceUploadRequestDTO) GetUploadedFilepath() string {
	return r.uploadedFilepath
}
func (r *ResourceUploadRequestDTO) SetUploadedFilepath(filepath string) {
	r.uploadedFilepath = filepath
}
func (r *ResourceUploadRequestDTO) GetUploadedFilesize() int64 {
	return r.uploadedFilesize
}
func (r *ResourceUploadRequestDTO) SetUploadedFilesize(filesize int64) {
	r.uploadedFilesize = filesize
}
func (r *ResourceUploadRequestDTO) GetUploadedFiletype() string {
	return r.uploadedFiletype
}
func (r *ResourceUploadRequestDTO) SetUploadedFiletype(filetype string) {
	r.uploadedFiletype = filetype
}
