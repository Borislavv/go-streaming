package uploader

import (
	"fmt"
	"github.com/Borislavv/video-streaming/internal/domain/dto"
	"github.com/Borislavv/video-streaming/internal/domain/errors"
	"github.com/Borislavv/video-streaming/internal/domain/logger"
	"github.com/Borislavv/video-streaming/internal/domain/service/storage"
	"github.com/Borislavv/video-streaming/internal/infrastructure/service/uploader/file"
	"io"
	"mime/multipart"
)

const partFilename = "resource"

// PartsUploader - is an file Uploader which use multipart.Part.
// In such case it takes more time but takes much less memory.
// Approximately, to upload a 50MB file you will need only 10MB of RAM.
type PartsUploader struct {
	logger                    logger.Logger
	storage                   storage.Storage
	filename                  file.NameComputer
	inMemoryFileSizeThreshold int64
}

func NewPartsUploader(
	logger logger.Logger,
	storage storage.Storage,
	filename file.NameComputer,
	inMemoryFileSizeThreshold int64,
) *PartsUploader {
	return &PartsUploader{
		logger:                    logger,
		filename:                  filename,
		storage:                   storage,
		inMemoryFileSizeThreshold: inMemoryFileSizeThreshold,
	}
}

func (u *PartsUploader) Upload(dto dto.UploadRequest) (err error) {
	part, err := u.getFilePart(dto)
	if err != nil {
		return u.logger.LogPropagate(err)
	}

	// TODO must be added filesize for check uniqueness
	computedFilename, err := u.filename.Get(
		part.FileName(),
		part.Header.Get("Content-Type"),
		part.Header.Get("Content-Disposition"),
	)

	// checking whether the being uploaded resource already exists
	has, err := u.storage.Has(computedFilename)
	if err != nil {
		return u.logger.LogPropagate(err)
	}
	if has { // if being uploading resource is already exists, then throw an error
		return u.logger.LogPropagate(errors.NewResourceAlreadyExistsError(part.FileName()))
	}

	// saving a file on disk and calculating new hashed name with full qualified path
	length, filename, filepath, err := u.storage.Store(computedFilename, part)
	if err != nil {
		return u.logger.LogPropagate(err)
	}

	// mutate request dto
	dto.SetOriginFilename(part.FileName())
	dto.SetUploadedFilename(filename)
	dto.SetUploadedFilepath(filepath)
	dto.SetUploadedFilesize(length)
	dto.SetUploadedFiletype(part.Header.Get("Content-Type"))

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
