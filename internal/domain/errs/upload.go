package errs

import (
	"fmt"
	"github.com/Borislavv/video-streaming/internal/infrastructure/logger"
	"net/http"
)

const (
	uploadErrType         = "application"
	publicUploadErrLevel  = logger.WarningLevel
	publicUploadErrStatus = http.StatusBadRequest
)

type ResourceAlreadyExistsError struct{ errored }

func NewResourceAlreadyExistsError(name string) *ResourceAlreadyExistsError {
	return &ResourceAlreadyExistsError{
		errored{
			ErrorMessage: fmt.Sprintf("the resource '%v' being loaded already exists", name),
			ErrorType:    uploadErrType,
			errorStatus:  publicUploadErrStatus,
			errorLevel:   publicUploadErrLevel,
		},
	}
}
