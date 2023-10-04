package uploader

import (
	"github.com/Borislavv/video-streaming/internal/domain/dto"
	"github.com/Borislavv/video-streaming/internal/domain/errors"
	"github.com/Borislavv/video-streaming/internal/domain/logger"
	"github.com/Borislavv/video-streaming/internal/domain/service/storage"
	"github.com/Borislavv/video-streaming/internal/infrastructure/service/uploader/file"
)

// NativeUploader is a service which represents functionality
// for uploader a full file from *http.Request into storage.
type NativeUploader struct {
	logger   logger.Logger
	storage  storage.Storage
	filename file.NameComputer
}

func NewNativeUploader(
	logger logger.Logger,
	storage storage.Storage,
	filename file.NameComputer,
) *NativeUploader {
	return &NativeUploader{
		logger:   logger,
		storage:  storage,
		filename: filename,
	}
}

// Upload method will be store a file on a disk and calculate a new hashed name. Request DTO mutation!
func (u *NativeUploader) Upload(req dto.UploadRequest) (err error) {
	computedFilename, err := u.filename.Get(
		req.GetPart().FileName(),
		req.GetPart().Header.Get("Content-Type"),
		req.GetPart().Header.Get("Content-Disposition"),
	)

	// checking whether the being uploaded resource already exists
	has, err := u.storage.Has(computedFilename)
	if err != nil {
		return u.logger.LogPropagate(err)
	}
	if has { // if being uploading resource is already exists, then throw an error
		return u.logger.LogPropagate(errors.NewResourceAlreadyExistsError(req.GetPart().FileName()))
	}

	// saving a file on disk and calculating new hashed name with full qualified path
	length, filename, filepath, err := u.storage.StoreConcurrently(computedFilename, req.GetPart())
	if err != nil {
		return u.logger.LogPropagate(err)
	}

	// mutate request dto
	req.SetUploadedFilename(filename)
	req.SetUploadedFilepath(filepath)
	req.SetUploadedFilesize(length)
	req.SetUploadedFiletype(req.GetPart().Header.Get("Content-Type"))

	return nil
}
