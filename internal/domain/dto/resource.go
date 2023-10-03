package dto

import (
	"mime/multipart"
)

type ResourceUploadRequestDTO struct {
	part             *multipart.Part
	contentLength    int64
	uploadedFilename string
	uploadedFilepath string
	uploadedFiletype string
	uploadedFilesize int64
}

func NewResourceUploadRequest(part *multipart.Part, contentLength int64) (dto *ResourceUploadRequestDTO) {
	return &ResourceUploadRequestDTO{part: part, contentLength: contentLength}
}

func (r *ResourceUploadRequestDTO) GetPart() *multipart.Part {
	return r.part
}
func (r *ResourceUploadRequestDTO) GetContentLength() int64 {
	return r.contentLength
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
