package dto

import (
	"mime/multipart"
	"runtime"
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
