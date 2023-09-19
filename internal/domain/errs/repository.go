package errs

import (
	"fmt"
	"github.com/Borislavv/video-streaming/internal/infrastructure/logger"
	"net/http"
)

const (
	repositoryType         = "application"
	publicRepositoryLevel  = logger.ErrorLevel
	publicRepositoryStatus = http.StatusNotFound
)

type NotFoundError struct{ errored }

func NewNotFoundError(entity string, by string) *NotFoundError {
	return &NotFoundError{
		errored{
			ErrorMessage: fmt.Sprintf("entity '%v' not found by given '%s'", entity, by),
			ErrorType:    repositoryType,
			errorLevel:   publicRepositoryLevel,
			errorStatus:  publicRepositoryStatus,
		},
	}
}

func IsNotFoundError(err error) bool {
	_, ok := err.(*NotFoundError)
	return ok
}
