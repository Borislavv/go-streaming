package errtype

import (
	"github.com/Borislavv/video-streaming/internal/infrastructure/service/logger"
	"net/http"
)

const (
	accessErrType         = "access"
	publicAccessErrLevel  = logger.InfoLevel
	publicAccessErrStatus = http.StatusBadRequest
)

type AccessDeniedError struct{ publicError }

func NewAccessDeniedError(msg string) *AccessDeniedError {
	baseMsg := "access denied"
	if msg != "" {
		msg = baseMsg + ": " + msg
	} else {
		msg = baseMsg
	}

	return &AccessDeniedError{
		publicError{
			errored{
				ErrorMessage: msg,
				ErrorType:    accessErrType,
				errorLevel:   publicAccessErrLevel,
				errorStatus:  publicAccessErrStatus,
			},
		},
	}
}
