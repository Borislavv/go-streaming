package errors

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

type InternalServerError struct{ publicError }

func NewInternalServerError(msg string) InternalServerError {
	if msg == "" {
		msg = internalServerErrorMessage
	}
	return InternalServerError{
		publicError{
			errored{
				ErrorMessage: msg,
				ErrorType:    internalServerErrorType,
				errorStatus:  internalServerErrorStatusCode,
				errorLevel:   internalServerErrorLevel,
			},
		},
	}
}
