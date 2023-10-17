package uploader

import (
	"github.com/Borislavv/video-streaming/internal/domain/dto"
	"github.com/Borislavv/video-streaming/internal/domain/errors"
	"github.com/Borislavv/video-streaming/internal/domain/logger"
	"github.com/Borislavv/video-streaming/internal/domain/service/storage"
	"github.com/Borislavv/video-streaming/internal/infrastructure/service/uploader/file"
	"io"
	"mime/multipart"
)

const MultipartPartUploadingType = "multipart_part"

// MultipartPartUploader - is a file UploadingStrategy which use multipart.Part.
// In such case it takes more time but takes much less memory.
// Approximately, to upload a 50MB file you will need only 10MB of RAM.
type MultipartPartUploader struct {
	logger      logger.Logger
	storage     storage.Storage
	filename    file.NameComputer
	maxFilesize int64
}

func NewPartsUploader(
	logger logger.Logger,
	storage storage.Storage,
	filename file.NameComputer,
) *MultipartPartUploader {
	return &MultipartPartUploader{
		logger:   logger,
		filename: filename,
		storage:  storage,
	}
}

func (u *MultipartPartUploader) Upload(reqDTO dto.UploadResourceRequest) (err error) {
	part, err := u.getFilePart(reqDTO)
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

	// mutate request reqDTO
	reqDTO.SetOriginFilename(part.FileName())
	reqDTO.SetUploadedFilename(filename)
	reqDTO.SetUploadedFilepath(filepath)
	reqDTO.SetUploadedFilesize(length)
	reqDTO.SetUploadedFiletype(part.Header.Get("Content-Type"))

	return nil
}

func (u *MultipartPartUploader) getFilePart(reqDTO dto.UploadResourceRequest) (part *multipart.Part, err error) {
	// extract the multipart form reader (handling the form as a stream)
	reader, err := reqDTO.GetRequest().MultipartReader()
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
		if part.FileName() != "" {
			return part, nil
		}
	}

	return nil, errors.NewFormDoesNotContainsUploadedFileError()
}
