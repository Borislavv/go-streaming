package _interface

import (
	"github.com/Borislavv/video-streaming/internal/domain/dto"
)

type Uploader interface {
	// Upload method will be store a file on a disk and calculate a new hashed name. Request DTO mutation!
	Upload(dto.UploadResourceRequest) (err error)
}