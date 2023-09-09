package errs

import (
	"github.com/Borislavv/video-streaming/internal/infrastructure/logger"
	"net/http"
)

const (
	internalServerErrorMessage    = "internal server error"
	internalServerErrorType       = "application"
	internalServerErrorStatusCode = http.StatusInternalServerError
	internalServerErrorLevel      = logger.ErrorLevel
)

type InternalServerError struct{ errored }

func NewInternalServerError() InternalServerError {
	return InternalServerError{
		errored{
			ErrorMessage: internalServerErrorMessage,
			ErrorType:    internalServerErrorType,
			errorStatus:  internalServerErrorStatusCode,
			errorLevel:   internalServerErrorLevel,
		},
	}
}
