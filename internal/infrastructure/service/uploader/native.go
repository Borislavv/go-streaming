package uploader

import (
	"github.com/Borislavv/video-streaming/internal/domain/dto"
	"github.com/Borislavv/video-streaming/internal/domain/errors"
	"github.com/Borislavv/video-streaming/internal/domain/logger"
	"github.com/Borislavv/video-streaming/internal/domain/service/storage"
	"github.com/Borislavv/video-streaming/internal/infrastructure/service/uploader/file"
)

const formFilename = "resource"

// NativeUploader is a service which represents functionality
// for uploader a full file from *http.Request into storage.
type NativeUploader struct {
	logger                    logger.Logger
	storage                   storage.Storage
	filename                  file.NameComputer
	inMemoryFileSizeThreshold int64
}

func NewNativeUploader(
	logger logger.Logger,
	storage storage.Storage,
	filename file.NameComputer,
	inMemoryFileSizeThreshold int64,
) *NativeUploader {
	return &NativeUploader{
		logger:                    logger,
		storage:                   storage,
		filename:                  filename,
		inMemoryFileSizeThreshold: inMemoryFileSizeThreshold,
	}
}

// Upload method will be store a file on a disk and calculate a new hashed name. Request DTO mutation!
func (u *NativeUploader) Upload(dto dto.UploadRequest) (err error) {
	// request will be parsed and stored in the memory if it is under the RAM threshold,
	// otherwise last parts of parsed file will be stored in the tmp files on the disk space
	if err = dto.GetRequest().ParseMultipartForm(u.inMemoryFileSizeThreshold); err != nil {
		return u.logger.LogPropagate(err)
	}

	// receiving a file and header from multipart/form-data
	// by requested filename which is stored in the `formFilename` const.
	formFile, header, err := dto.GetRequest().FormFile(formFilename)
	if err != nil {
		return u.logger.LogPropagate(err)
	}
	defer func() { _ = formFile.Close() }()

	computedFilename, err := u.filename.Get(
		dto.GetPart().FileName(),
		dto.GetPart().Header.Get("Content-Type"),
		dto.GetPart().Header.Get("Content-Disposition"),
	)

	// checking whether the being uploaded resource already exists
	has, err := u.storage.Has(computedFilename)
	if err != nil {
		return u.logger.LogPropagate(err)
	}
	if has { // if being uploading resource is already exists, then throw an error
		return u.logger.LogPropagate(errors.NewResourceAlreadyExistsError(dto.GetPart().FileName()))
	}

	// saving a file on disk and calculating new hashed name with full qualified path
	length, filename, filepath, err := u.storage.Store(computedFilename, dto.GetPart())
	if err != nil {
		return u.logger.LogPropagate(err)
	}

	// mutate request dto
	dto.SetUploadedFilename(filename)
	dto.SetUploadedFilepath(filepath)
	dto.SetUploadedFilesize(length)
	dto.SetUploadedFiletype(dto.GetPart().Header.Get("Content-Type"))

	return nil
}
