package errs

import (
	"fmt"
	"net/http"
)

const ValidationType = "validation"

type FieldCannotBeEmptyError struct {
	Message string
	Type    string
}

func NewFieldCannotBeEmptyError(field string) *FieldCannotBeEmptyError {
	return &FieldCannotBeEmptyError{
		Message: fmt.Sprintf("field '%s' must not be empty or omitted", field),
		Type:    ValidationType,
	}
}

func (e *FieldCannotBeEmptyError) Error() string {
	return e.Message
}

func (e *FieldCannotBeEmptyError) Status() int {
	return http.StatusBadRequest
}
