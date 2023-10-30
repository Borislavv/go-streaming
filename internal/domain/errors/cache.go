package errors

import (
	"fmt"
	"github.com/Borislavv/video-streaming/internal/infrastructure/service/logger"
	"net/http"
)

const (
	cacheType                          = "cache"
	cacheInternalServerErrorStatusCode = http.StatusInternalServerError
	cacheInternalServerErrorLevel      = logger.CriticalLevel
)

type CachedDataTypeWasNotMatchedError struct{ internalError }

func NewCachedDataTypeWasNotMatchedError(key string) *CachedDataTypeWasNotMatchedError {
	return &CachedDataTypeWasNotMatchedError{
		internalError{
			errored{
				ErrorMessage: fmt.Sprintf("cached data type was not matched by key '%v'", key),
				ErrorType:    cacheType,
				errorLevel:   cacheInternalServerErrorLevel,
				errorStatus:  cacheInternalServerErrorStatusCode,
			},
		},
	}
}
