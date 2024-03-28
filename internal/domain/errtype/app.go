package errtype

import (
	"github.com/Borislavv/video-streaming/internal/infrastructure/service/logger"
	"net/http"
)

const (
	applicationType               = "application"
	internalServerErrorMessage    = "internal server error"
	internalServerErrorStatusCode = http.StatusInternalServerError
	internalServerErrorLevel      = logger.ErrorLevel
)

type InternalServerError struct{ publicError }

func NewInternalServerError() InternalServerError {
	return InternalServerError{
		publicError{
			errored{
				ErrorType:    applicationType,
				ErrorMessage: internalServerErrorMessage,
				errorStatus:  internalServerErrorStatusCode,
				errorLevel:   internalServerErrorLevel,
			},
		},
	}
}
