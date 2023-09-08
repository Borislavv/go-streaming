package errs

import (
	"fmt"
	"net/http"
	"strings"
)

const (
	ValidationType   = "validation"
	ValidationStatus = http.StatusBadRequest
)

type FieldCannotBeEmptyError Error

func NewFieldCannotBeEmptyError(field string) *FieldCannotBeEmptyError {
	return &FieldCannotBeEmptyError{
		Message: fmt.Sprintf("field '%v' must not be empty or omitted", field),
		Type:    ValidationType,
	}
}

func (e *FieldCannotBeEmptyError) Error() string {
	return e.Message
}

func (e *FieldCannotBeEmptyError) Status() int {
	return ValidationStatus
}

type AtLeastOneFieldMustBeDefinedError Error

func NewAtLeastOneFieldMustBeDefinedError(fields ...string) *AtLeastOneFieldMustBeDefinedError {
	return &AtLeastOneFieldMustBeDefinedError{
		Message: fmt.Sprintf("at least one of followed fields [%v] must be defined", strings.Join(fields, ", ")),
		Type:    ValidationType,
	}
}

func (e *AtLeastOneFieldMustBeDefinedError) Error() string {
	return e.Message
}

func (e *AtLeastOneFieldMustBeDefinedError) Status() int {
	return ValidationStatus
}

type FieldLengthMustBeMoreOrLessError Error

func NewFieldLengthMustBeMoreOrLessError(field string, isMustBeMore bool, length int) *FieldLengthMustBeMoreOrLessError {
	qualifier := "less"
	if isMustBeMore {
		qualifier = "more"
	}
	return &FieldLengthMustBeMoreOrLessError{
		Message: fmt.Sprintf("length of the field '%v' must be %v than %d", field, qualifier, length),
		Type:    ValidationType,
	}
}

func (e *FieldLengthMustBeMoreOrLessError) Error() string {
	return e.Message
}

func (e *FieldLengthMustBeMoreOrLessError) Status() int {
	return ValidationStatus
}

type UniquenessCheckFailedError Error

func NewUniquenessCheckFailedError(field string) *UniquenessCheckFailedError {
	return &UniquenessCheckFailedError{
		Message: fmt.Sprintf("uniqueness check filed due to duplicated '%v'", field),
		Type:    ValidationType,
	}
}

func (e *UniquenessCheckFailedError) Error() string {
	return e.Message
}

func (e *UniquenessCheckFailedError) Status() int {
	return ValidationStatus
}
