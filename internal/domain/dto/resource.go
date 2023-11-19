package dto

import (
	"github.com/Borislavv/video-streaming/internal/domain/vo"
	"net/http"
)

type ResourceUploadRequestDTO struct {
	userID           vo.ID
	request          *http.Request
	originFilename   string
	contentLength    int64
	uploadedFilename string
	uploadedFilepath string
	uploadedFiletype string
	uploadedFilesize int64
}

func NewResourceUploadRequest(r *http.Request) (dto *ResourceUploadRequestDTO) {
	return &ResourceUploadRequestDTO{request: r}
}

func (r *ResourceUploadRequestDTO) GetUserID() vo.ID {
	return r.userID
}
func (r *ResourceUploadRequestDTO) SetUserID(id vo.ID) {
	r.userID = id
}
func (r *ResourceUploadRequestDTO) GetRequest() *http.Request {
	return r.request
}
func (r *ResourceUploadRequestDTO) GetOriginFilename() string {
	return r.originFilename
}
func (r *ResourceUploadRequestDTO) SetOriginFilename(filename string) {
	r.originFilename = filename
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

type ResourceGetRequestDTO struct {
	ID     vo.ID `json:"id"`
	UserID vo.ID
}

func NewResourceGetRequestDTO(id vo.ID, userID vo.ID) *ResourceGetRequestDTO {
	return &ResourceGetRequestDTO{
		ID:     id,
		UserID: userID,
	}
}
func (req *ResourceGetRequestDTO) GetID() vo.ID {
	return req.ID
}
func (req *ResourceGetRequestDTO) GetUserID() vo.ID {
	return req.UserID
}

type ResourceDeleteRequestDTO struct {
	ID     vo.ID `json:"id"`
	UserID vo.ID
}

func NewResourceDeleteRequestDTO(id vo.ID, userID vo.ID) *ResourceDeleteRequestDTO {
	return &ResourceDeleteRequestDTO{
		ID:     id,
		UserID: userID,
	}
}
func (req *ResourceDeleteRequestDTO) GetID() vo.ID {
	return req.ID
}
func (req *ResourceDeleteRequestDTO) GetUserID() vo.ID {
	return req.UserID
}
