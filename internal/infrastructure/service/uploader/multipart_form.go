package uploader

import (
	"github.com/Borislavv/video-streaming/internal/domain/dto"
	"github.com/Borislavv/video-streaming/internal/domain/errors"
	"github.com/Borislavv/video-streaming/internal/domain/logger"
	"github.com/Borislavv/video-streaming/internal/domain/service/storage"
	"github.com/Borislavv/video-streaming/internal/infrastructure/service/uploader/file"
)

// MultipartFormUploader is a service which represents functionality
// for uploader a full file from *http.Request into storage.
// This approach of uploading takes a much more RAM but works more fast than MultipartPartUploader.
// If you care of performance, you need use this approach, but take care of using RAM and set up
// the appropriate value of 'inMemoryFileSizeThreshold' through env. configuration.
type MultipartFormUploader struct {
	logger                    logger.Logger
	storage                   storage.Storage
	filename                  file.NameComputer
	formFilename              string
	maxFilesize               int64
	inMemoryFileSizeThreshold int64
}

func NewNativeUploader(
	logger logger.Logger,
	storage storage.Storage,
	filename file.NameComputer,
	formFilename string,
	inMemoryFileSizeThreshold int64,
) *MultipartFormUploader {
	return &MultipartFormUploader{
		logger:                    logger,
		storage:                   storage,
		filename:                  filename,
		formFilename:              formFilename,
		inMemoryFileSizeThreshold: inMemoryFileSizeThreshold,
	}
}

// Upload method will be store a file on a disk and calculate a new hashed name. Request DTO mutation!
func (u *MultipartFormUploader) Upload(dto dto.UploadRequest) (err error) {
	// request will be parsed and stored in the memory if it is under the RAM threshold,
	// otherwise last parts of parsed file will be stored in the tmp files on the disk space
	if err = dto.GetRequest().ParseMultipartForm(u.inMemoryFileSizeThreshold); err != nil {
		return u.logger.LogPropagate(err)
	}

	// receiving a file and header from multipart/form-data
	// by requested filename which is stored in the `formFilename` const.
	formFile, header, err := dto.GetRequest().FormFile(u.formFilename)
	if err != nil {
		return u.logger.LogPropagate(err)
	}
	defer func() { _ = formFile.Close() }()

	// TODO must be added filesize for check uniqueness
	computedFilename, err := u.filename.Get(
		header.Filename,
		header.Header.Get("Content-Type"),
		header.Header.Get("Content-Disposition"),
	)

	// checking whether the being uploaded resource already exists
	has, err := u.storage.Has(computedFilename)
	if err != nil {
		return u.logger.LogPropagate(err)
	}
	if has { // if being uploading resource is already exists, then throw an error
		return u.logger.LogPropagate(errors.NewResourceAlreadyExistsError(header.Filename))
	}

	// saving a file on disk and calculating new hashed name with full qualified path
	length, filename, filepath, err := u.storage.Store(computedFilename, formFile)
	if err != nil {
		return u.logger.LogPropagate(err)
	}

	// mutate request dto
	dto.SetOriginFilename(header.Filename)
	dto.SetUploadedFilename(filename)
	dto.SetUploadedFilepath(filepath)
	dto.SetUploadedFilesize(length)
	dto.SetUploadedFiletype(header.Header.Get("Content-Type"))

	return nil
}
