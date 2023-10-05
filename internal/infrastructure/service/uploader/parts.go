package uploader

import (
	"fmt"
	"github.com/Borislavv/video-streaming/internal/domain/dto"
	"github.com/Borislavv/video-streaming/internal/domain/errors"
	"github.com/Borislavv/video-streaming/internal/domain/logger"
	"io"
	"mime/multipart"
)

const partFilename = "resource"

// PartsUploader - is an file Uploader which use multipart.Part.
// In such case it takes more time but takes much less memory.
// Approximately, to upload a 50MB file you will need only 10MB of RAM.
type PartsUploader struct {
	logger logger.Logger
}

func NewPartsUploader(logger logger.Logger) *PartsUploader {
	return &PartsUploader{logger: logger}
}

func (u *PartsUploader) Upload(dto dto.UploadRequest) (err error) {
	// todo must be implemented
	return nil
}

func (u *PartsUploader) getFilePart(dto dto.UploadRequest) (part *multipart.Part, err error) {
	// extract the multipart form reader (handling the form as a stream)
	reader, err := dto.GetRequest().MultipartReader()
	if err != nil {
		return nil, u.logger.LogPropagate(err)
	}

	for { // find the part of the form with the target file
		part, err = reader.NextPart()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, u.logger.LogPropagate(err)
		}

		// check the form part is th target file field
		if part.FileName() == partFilename {
			return part, nil
		}
	}

	return nil, errors.NewInvalidUploadedFileError(
		fmt.Sprintf("form does not contains the target file field '%v'", partFilename),
	)
}
