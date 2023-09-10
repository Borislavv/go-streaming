package uploader

import (
	"fmt"
	"github.com/Borislavv/video-streaming/internal/domain/errs"
	"github.com/Borislavv/video-streaming/internal/domain/service"
	"net/http"
)

const (
	maxFileSize  = 100 << 20 // 100mb.
	formFilename = "resource"
)

// NativeUploader is a service which represents functionality
// for uploader a full file from *http.Request into storage.
type NativeUploader struct {
	logger  service.Logger
	storage service.Storage
}

func NewNativeUploader(
	logger service.Logger,
	storage service.Storage,
) *NativeUploader {
	return &NativeUploader{
		logger:  logger,
		storage: storage,
	}
}

func (u *NativeUploader) Upload(r *http.Request) (resourceId string, e error) {
	// request parsing into memory if it's under the threshold or in tmp files
	// the max value of allowed RAM memory for each file is stored in the `maxFileSize` const.
	if err := r.ParseMultipartForm(maxFileSize); err != nil {
		return "", u.logger.LogPropagate(err)
	}

	// receiving a file and header from multipart/form-data
	// by requested filename which is stored in the `formFilename` const.
	file, header, err := r.FormFile(formFilename)
	if err != nil {
		return "", u.logger.LogPropagate(err)
	}
	defer file.Close()

	// logging the uploaded file info
	u.logger.Info(
		fmt.Sprintf(
			"\n\t\t\tUploaded:\n"+
				"\t\t\t\tFilename: %v\n"+
				"\t\t\t\tFilesize: %d\n"+
				"\t\t\t\tMIME type: %v",
			header.Filename,
			header.Size,
			header.Header,
		),
	)

	// checking whether the being uploaded resource already exists
	has, err := u.storage.Has(header)
	if err != nil {
		return "", u.logger.LogPropagate(err)
	}
	if has { // if being uploading resource is already exists, then throw an error
		return "", u.logger.LogPropagate(errs.NewResourceAlreadyExistsError(header.Filename))
	}

	id, err := u.storage.Store(file, header)
	if err != nil {
		return "", err
	}

	return id, nil
}
