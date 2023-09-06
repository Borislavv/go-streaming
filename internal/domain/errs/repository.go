package errs

import (
	"fmt"
	"net/http"
)

const RepositoryType = "repository"

type NotFoundError Error

func NewNotFoundError(entity string) *NotFoundError {
	return &NotFoundError{
		Message: fmt.Sprintf("%s not found by given id", entity),
		Type:    RepositoryType,
	}
}

func (e *NotFoundError) Error() string {
	return e.Message
}

func (e *NotFoundError) Status() int {
	return http.StatusInternalServerError
}
