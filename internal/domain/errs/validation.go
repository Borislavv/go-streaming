package errs

import (
	"fmt"
	"net/http"
)

const ValidationType = "validation"

type FieldCannotBeEmptyError struct {
	Message string `json:"message"`
	Type    string `json:"type"`
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

type UniquenessCheckFailedError struct {
	Message string `json:"message"`
	Type    string `json:"type"`
}

func NewUniquenessCheckFailedError(field string) *UniquenessCheckFailedError {
	return &UniquenessCheckFailedError{
		Message: fmt.Sprintf("uniqueness check filed due to duplicated '%s'", field),
		Type:    ValidationType,
	}
}

func (e *UniquenessCheckFailedError) Error() string {
	return e.Message
}

func (e *UniquenessCheckFailedError) Status() int {
	return http.StatusBadRequest
}
