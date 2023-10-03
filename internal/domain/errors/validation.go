package errors

import (
	"fmt"
	"github.com/Borislavv/video-streaming/internal/infrastructure/logger"
	"net/http"
	"strings"
)

const (
	validationType = "validation"

	publicValidationStatus = http.StatusBadRequest
	publicValidationLevel  = logger.WarningLevel

	internalValidationStatus = http.StatusInternalServerError
	internalValidationLevel  = logger.CriticalLevel
)

type InternalValidationError struct{ internalError }

func NewInternalValidationError(message string) *InternalValidationError {
	return &InternalValidationError{
		internalError{
			errored{
				ErrorMessage: message,
				ErrorType:    validationType,
				errorLevel:   internalValidationLevel,
				errorStatus:  internalValidationStatus,
			},
		},
	}
}

type FieldCannotBeEmptyError struct{ publicError }

func NewFieldCannotBeEmptyError(field string) *FieldCannotBeEmptyError {
	return &FieldCannotBeEmptyError{
		publicError{
			errored{
				ErrorMessage: fmt.Sprintf("field '%v' must not be empty or omitted", field),
				ErrorType:    validationType,
				errorLevel:   publicValidationLevel,
				errorStatus:  publicValidationStatus,
			},
		},
	}
}

type AtLeastOneFieldMustBeDefinedError struct{ publicError }

func NewAtLeastOneFieldMustBeDefinedError(fields ...string) *AtLeastOneFieldMustBeDefinedError {
	return &AtLeastOneFieldMustBeDefinedError{
		publicError{
			errored{
				ErrorMessage: fmt.Sprintf("at least one of followed fields [%v] must be defined", strings.Join(fields, ", ")),
				ErrorType:    validationType,
				errorLevel:   publicValidationLevel,
				errorStatus:  publicValidationStatus,
			},
		},
	}
}

type FieldLengthMustBeMoreOrLessError struct{ publicError }

func NewFieldLengthMustBeMoreOrLessError(field string, isMustBeMore bool, length int) *FieldLengthMustBeMoreOrLessError {
	qualifier := "less"
	if isMustBeMore {
		qualifier = "more"
	}
	return &FieldLengthMustBeMoreOrLessError{
		publicError{
			errored{
				ErrorMessage: fmt.Sprintf("length of the field '%v' must be %v than %d", field, qualifier, length),
				ErrorType:    validationType,
				errorLevel:   publicValidationLevel,
				errorStatus:  publicValidationStatus,
			},
		},
	}
}

type UniquenessCheckFailedError struct{ publicError }

func NewUniquenessCheckFailedError(fields ...string) *UniquenessCheckFailedError {
	return &UniquenessCheckFailedError{
		publicError{
			errored{
				ErrorMessage: fmt.Sprintf(
					"uniqueness check failed due to duplicated key '%v'", strings.Join(fields, ", "),
				),
				ErrorType:   validationType,
				errorLevel:  publicValidationLevel,
				errorStatus: publicValidationStatus,
			},
		},
	}
}

type InvalidUploadedFileError struct{ publicError }

func NewInvalidUploadedFileError(message string) *InvalidUploadedFileError {
	if message == "" {
		message = "given file is not valid"
	}

	return &InvalidUploadedFileError{
		publicError{
			errored{
				ErrorMessage: message,
				ErrorType:    validationType,
				errorLevel:   publicValidationLevel,
				errorStatus:  publicValidationStatus,
			},
		},
	}
}

type FormDoesNotContainsUploadedFileError struct{ publicError }

func NewFormDoesNotContainsUploadedFileError() *FormDoesNotContainsUploadedFileError {
	return &FormDoesNotContainsUploadedFileError{
		publicError{
			errored{
				ErrorMessage: "form does not contains an uploading file",
				ErrorType:    validationType,
				errorLevel:   publicValidationLevel,
				errorStatus:  publicValidationStatus,
			},
		},
	}
}

type TimeParsingValidationError struct{ publicError }

func NewTimeParsingValidationError(date string) *TimeParsingValidationError {
	return &TimeParsingValidationError{
		publicError{
			errored{
				ErrorMessage: fmt.Sprintf(
					"date string '%v' has wrong format. "+
						"date + time + timezone - corrent format: '2006-01-02T15:04:05-07:00', "+
						"date + time - corrent format: '2006-01-02T15:04:05', "+
						"date - corrent format: '2006-01-02'",
					date,
				),
				ErrorType:   validationType,
				errorLevel:  publicValidationLevel,
				errorStatus: publicValidationStatus,
			},
		},
	}
}
