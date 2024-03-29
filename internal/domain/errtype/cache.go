package errtype

import (
	"fmt"
	"github.com/Borislavv/video-streaming/internal/infrastructure/service/logger"
	"net/http"
	"reflect"
)

const (
	cacheType                          = "cache"
	cacheInternalServerErrorStatusCode = http.StatusInternalServerError
	cacheInternalServerErrorLevel      = logger.CriticalLevel
)

type CachedDataTypeWasNotMatchedError struct{ internalError }

func NewCachedDataTypeWasNotMatchedError(
	key string, expectedType reflect.Type, fetchedType reflect.Type,
) *CachedDataTypeWasNotMatchedError {
	return &CachedDataTypeWasNotMatchedError{
		internalError{
			errored{
				ErrorMessage: fmt.Sprintf(
					"cached type of data was not matched by key '%v', expected: '%v', fetched: '%v'",
					key, expectedType.String(), fetchedType.String(),
				),
				ErrorType:   cacheType,
				errorLevel:  cacheInternalServerErrorLevel,
				errorStatus: cacheInternalServerErrorStatusCode,
			},
		},
	}
}
