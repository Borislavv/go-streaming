package dto

import (
	"mime/multipart"
	"runtime"
)

type ResourceUploadRequestDTO struct {
	file             multipart.File
	header           *multipart.FileHeader
	uploadedFilename string
	uploadedFilepath string
}

func NewResourceUploadRequest(
	file multipart.File,
	header *multipart.FileHeader,
) (dto *ResourceUploadRequestDTO) {
	resourceDTO := &ResourceUploadRequestDTO{
		file:   file,
		header: header,
	}

	// destructor implementation for deferred close file
	runtime.SetFinalizer(resourceDTO, func(dto *ResourceUploadRequestDTO) {
		_ = dto.file.Close()
	})

	return resourceDTO
}

func (r *ResourceUploadRequestDTO) GetFile() multipart.File {
	return r.file
}
func (r *ResourceUploadRequestDTO) GetHeader() *multipart.FileHeader {
	return r.header
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
