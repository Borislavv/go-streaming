package errtype

import (
	"fmt"
	"github.com/Borislavv/video-streaming/internal/infrastructure/service/logger"
	"net/http"
	"reflect"
)

const (
	serviceContainerType                          = "service-container"
	serviceContainerInternalServerErrorLevel      = logger.CriticalLevel
	serviceContainerInternalServerErrorStatusCode = http.StatusInternalServerError
)

type ServiceWasNotFoundIntoContainerError struct{ internalError }

func NewServiceWasNotFoundIntoContainerError(key reflect.Type) *ServiceWasNotFoundIntoContainerError {
	return &ServiceWasNotFoundIntoContainerError{
		internalError{
			errored{
				ErrorMessage: fmt.Sprintf("service '%v' was not found into the service container", key.String()),
				ErrorType:    serviceContainerType,
				errorLevel:   serviceContainerInternalServerErrorLevel,
				errorStatus:  serviceContainerInternalServerErrorStatusCode,
			},
		},
	}
}

type TypesMismatchedServiceContainerError struct{ internalError }

func NewTypesMismatchedServiceContainerError(
	given reflect.Type, expected reflect.Type,
) *TypesMismatchedServiceContainerError {
	return &TypesMismatchedServiceContainerError{
		internalError{
			errored{
				ErrorMessage: fmt.Sprintf("types mismatched error, service container "+
					"returned type of '%v' but expected '%v'", given.String(), expected.String()),
				ErrorType:   serviceContainerType,
				errorLevel:  serviceContainerInternalServerErrorLevel,
				errorStatus: serviceContainerInternalServerErrorStatusCode,
			},
		},
	}
}
