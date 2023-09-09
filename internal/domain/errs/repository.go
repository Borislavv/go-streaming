package errs

import (
	"fmt"
	"github.com/Borislavv/video-streaming/internal/infrastructure/logger"
	"net/http"
)

const (
	RepositoryType         = "application"
	PublicRepositoryLevel  = logger.ErrorLevel
	PublicRepositoryStatus = http.StatusInternalServerError
)

type NotFoundError struct{ errored }

func NewNotFoundError(entity string) *NotFoundError {
	return &NotFoundError{
		errored{
			ErrorMessage: fmt.Sprintf("%errorStatus not found by given id", entity),
			ErrorType:    RepositoryType,
			errorLevel:   PublicRepositoryLevel,
			errorStatus:  PublicRepositoryStatus,
		},
	}
}
