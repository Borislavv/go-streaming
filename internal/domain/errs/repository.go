package errs

import (
	"fmt"
	"github.com/Borislavv/video-streaming/internal/infrastructure/logger"
	"net/http"
)

const (
	repositoryType         = "application"
	publicRepositoryLevel  = logger.ErrorLevel
	publicRepositoryStatus = http.StatusInternalServerError
)

type NotFoundError struct{ errored }

func NewNotFoundError(entity string) *NotFoundError {
	return &NotFoundError{
		errored{
			ErrorMessage: fmt.Sprintf("%errorStatus not found by given id", entity),
			ErrorType:    repositoryType,
			errorLevel:   publicRepositoryLevel,
			errorStatus:  publicRepositoryStatus,
		},
	}
}
