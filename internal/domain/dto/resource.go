package dto

import (
	"mime/multipart"
	"runtime"
	"time"
)

type ResourceUploadRequestDto struct {
	file             multipart.File
	header           *multipart.FileHeader
	uploadedFilename string
	uploadedFilepath string
}

func NewResourceUploadRequest(
	file multipart.File,
	header *multipart.FileHeader,
) (dto *ResourceUploadRequestDto) {
	resourceDto := &ResourceUploadRequestDto{
		file:   file,
		header: header,
	}

	// destructor implementation for deferred close file
	runtime.SetFinalizer(resourceDto, func(dto *ResourceUploadRequestDto) {
		_ = dto.file.Close()
	})

	return resourceDto
}

func (r *ResourceUploadRequestDto) GetFile() multipart.File {
	return r.file
}
func (r *ResourceUploadRequestDto) GetHeader() *multipart.FileHeader {
	return r.header
}
func (r *ResourceUploadRequestDto) GetUploadedFilename() string {
	return r.uploadedFilename
}
func (r *ResourceUploadRequestDto) SetUploadedFilename(filename string) {
	r.uploadedFilename = filename
}
func (r *ResourceUploadRequestDto) GetUploadedFilepath() string {
	return r.uploadedFilepath
}
func (r *ResourceUploadRequestDto) SetUploadedFilepath(filepath string) {
	r.uploadedFilepath = filepath
}

type ResourceListRequestDto struct {
	CreatedAt time.Time `json:"createdAt" format:"2006-01-02T15:04:05Z07:00"`
	From      time.Time `json:"from" format:"2006-01-02T15:04:05Z07:00"`
	To        time.Time `json:"to" format:"2006-01-02T15:04:05Z07:00"`
	PaginationRequestDto
}

func (r *ResourceListRequestDto) GetCreatedAt() time.Time {
	return r.CreatedAt
}
func (r *ResourceListRequestDto) GetFrom() time.Time {
	return r.From
}
func (r *ResourceListRequestDto) GetTo() time.Time {
	return r.To
}
