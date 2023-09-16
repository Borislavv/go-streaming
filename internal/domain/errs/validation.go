package errs

import (
	"fmt"
	"github.com/Borislavv/video-streaming/internal/infrastructure/logger"
	"net/http"
	"strings"
)

const (
	validationType   = "validation"
	validationStatus = http.StatusBadRequest
	validationLevel  = logger.WarningLevel
)

type ValidationError struct{ errored }

func NewValidationError(message string) *ValidationError {
	return &ValidationError{
		errored{
			ErrorMessage: message,
			ErrorType:    validationType,
			errorLevel:   validationLevel,
			errorStatus:  validationStatus,
		},
	}
}

type FieldCannotBeEmptyError struct{ errored }

func NewFieldCannotBeEmptyError(field string) *FieldCannotBeEmptyError {
	return &FieldCannotBeEmptyError{
		errored{
			ErrorMessage: fmt.Sprintf("field '%v' must not be empty or omitted", field),
			ErrorType:    validationType,
			errorLevel:   validationLevel,
			errorStatus:  validationStatus,
		},
	}
}

type AtLeastOneFieldMustBeDefinedError struct{ errored }

func NewAtLeastOneFieldMustBeDefinedError(fields ...string) *AtLeastOneFieldMustBeDefinedError {
	return &AtLeastOneFieldMustBeDefinedError{
		errored{
			ErrorMessage: fmt.Sprintf("at least one of followed fields [%v] must be defined", strings.Join(fields, ", ")),
			ErrorType:    validationType,
			errorLevel:   validationLevel,
			errorStatus:  validationStatus,
		},
	}
}

type FieldLengthMustBeMoreOrLessError struct{ errored }

func NewFieldLengthMustBeMoreOrLessError(field string, isMustBeMore bool, length int) *FieldLengthMustBeMoreOrLessError {
	qualifier := "less"
	if isMustBeMore {
		qualifier = "more"
	}
	return &FieldLengthMustBeMoreOrLessError{
		errored{
			ErrorMessage: fmt.Sprintf("length of the field '%v' must be %v than %d", field, qualifier, length),
			ErrorType:    validationType,
			errorLevel:   validationLevel,
			errorStatus:  validationStatus,
		},
	}
}

type UniquenessCheckFailedError struct{ errored }

func NewUniquenessCheckFailedError(fields ...string) *UniquenessCheckFailedError {
	return &UniquenessCheckFailedError{
		errored{
			ErrorMessage: fmt.Sprintf(
				"uniqueness check failed due to duplicated key '%v'", strings.Join(fields, ", "),
			),
			ErrorType:   validationType,
			errorLevel:  validationLevel,
			errorStatus: validationStatus,
		},
	}
}

type InvalidUploadedFileError struct{ errored }

func NewInvalidUploadedFileError(filename string) *InvalidUploadedFileError {
	return &InvalidUploadedFileError{
		errored{
			ErrorMessage: fmt.Sprintf("file '%v' has a zero size", filename),
			ErrorType:    validationType,
			errorLevel:   validationLevel,
			errorStatus:  validationStatus,
		},
	}
}

type TimeParsingValidationError struct{ errored }

func NewTimeParsingValidationError(date string) *TimeParsingValidationError {
	return &TimeParsingValidationError{
		errored{
			ErrorMessage: fmt.Sprintf(
				"date string '%v' has wrong format. "+
					"date + time + timezone - corrent format: '2006-01-02T15:04:05-07:00', "+
					"date + time - corrent format: '2006-01-02T15:04:05', "+
					"date - corrent format: '2006-01-02'",
				date,
			),
			ErrorType:   validationType,
			errorLevel:  validationLevel,
			errorStatus: validationStatus,
		},
	}
}
