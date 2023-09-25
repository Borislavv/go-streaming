package uploader

import (
	"fmt"
	"github.com/Borislavv/video-streaming/internal/domain/dto"
	"github.com/Borislavv/video-streaming/internal/domain/errors"
	"github.com/Borislavv/video-streaming/internal/domain/logger"
	"github.com/Borislavv/video-streaming/internal/domain/service/storage"
)

// NativeUploader is a service which represents functionality
// for uploader a full file from *http.Request into storage.
type NativeUploader struct {
	logger  logger.Logger
	storage storage.Storage
}

func NewNativeUploader(
	logger logger.Logger,
	storage storage.Storage,
) *NativeUploader {
	return &NativeUploader{
		logger:  logger,
		storage: storage,
	}
}

// Upload method will be store a file on a disk and calculate a new hashed name. Request DTO mutation!
func (u *NativeUploader) Upload(req dto.UploadRequest) (err error) {
	// logging the uploaded file info
	u.logger.Info(
		fmt.Sprintf(
			"\n\t\t\tUploaded:\n"+
				"\t\t\t\tFilename: %v\n"+
				"\t\t\t\tFilesize: %d\n"+
				"\t\t\t\tMIME type: %v",
			req.GetHeader().Filename,
			req.GetHeader().Size,
			req.GetHeader().Header,
		),
	)

	// checking whether the being uploaded resource already exists
	has, err := u.storage.Has(req.GetHeader())
	if err != nil {
		return u.logger.LogPropagate(err)
	}
	if has { // if being uploading resource is already exists, then throw an error
		return u.logger.LogPropagate(errors.NewResourceAlreadyExistsError(req.GetHeader().Filename))
	}

	// saving a file on disk and calculating new hashed name with full qualified path
	filename, filepath, err := u.storage.Store(req.GetFile(), req.GetHeader())
	if err != nil {
		return u.logger.LogPropagate(err)
	}

	// mutate request dto
	req.SetUploadedFilename(filename)
	req.SetUploadedFilepath(filepath)

	return nil
}
