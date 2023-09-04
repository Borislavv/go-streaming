package errs

import "fmt"

const RepositoryType = "repository"

type NotFoundError struct {
	Message string
}

func NewNotFoundError(entity string) *NotFoundError {
	return &NotFoundError{
		Message: fmt.Sprintf("%s not found by given id", entity),
	}
}

func (e *NotFoundError) Error() string {
	return e.Message
}

func (e *NotFoundError) Public() bool {
	return true
}
