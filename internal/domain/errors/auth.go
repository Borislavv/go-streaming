package errors

import (
	"fmt"
	"github.com/Borislavv/video-streaming/internal/infrastructure/logger"
	"net/http"
)

const (
	authErrType           = "authorization"
	internalAuthErrLevel  = logger.CriticalLevel
	internalAuthErrStatus = http.StatusInternalServerError
	publicAuthErrLevel    = logger.InfoLevel
	publicAuthErrStatus   = http.StatusBadRequest
)

type AuthFailedError struct{ PublicError }

func NewAuthFailedError(msg string) *AuthFailedError {
	baseMsg := "authorization failed"
	if msg != "" {
		msg = fmt.Sprintf("%v: %v", baseMsg, msg)
	} else {
		msg = baseMsg
	}

	return &AuthFailedError{
		publicError{
			errored{
				ErrorMessage: msg,
				ErrorType:    authErrType,
				errorStatus:  publicAuthErrStatus,
				errorLevel:   publicAuthErrLevel,
			},
		},
	}
}

type TokenAlgoWasNotMatchedError struct{ internalError }

func NewTokenAlgoWasNotMatchedError() *TokenAlgoWasNotMatchedError {
	return &TokenAlgoWasNotMatchedError{
		internalError{
			errored{
				ErrorMessage: "token algo was not matched",
				ErrorType:    authErrType,
				errorStatus:  internalAuthErrStatus,
				errorLevel:   internalAuthErrLevel,
			},
		},
	}
}

type TokenUnexpectedSigningMethodError struct{ internalError }

func NewTokenUnexpectedSigningMethodError(algo interface{}) *TokenUnexpectedSigningMethodError {
	return &TokenUnexpectedSigningMethodError{
		internalError{
			errored{
				ErrorMessage: fmt.Sprintf("unexpected signing algo '%v'", algo),
				ErrorType:    authErrType,
				errorStatus:  internalAuthErrStatus,
				errorLevel:   internalAuthErrLevel,
			},
		},
	}
}

type TokenSubjectPayloadWasNotMatchedError struct{ internalError }

func NewTokenSubjectPayloadWasNotMatchedError(expected string, given string) *TokenSubjectPayloadWasNotMatchedError {
	return &TokenSubjectPayloadWasNotMatchedError{
		internalError{
			errored{
				ErrorMessage: fmt.Sprintf(
					"subject payload was not matched, expected: %v, given: %v", expected, given,
				),
				ErrorType:   authErrType,
				errorStatus: internalAuthErrStatus,
				errorLevel:  internalAuthErrLevel,
			},
		},
	}
}

type TokenInvalidError struct{ internalError }

func NewTokenInvalidError(token string) *TokenInvalidError {
	return &TokenInvalidError{
		internalError{
			errored{
				ErrorMessage: fmt.Sprintf("given token '%v' is not valid", token),
				ErrorType:    authErrType,
				errorStatus:  internalAuthErrStatus,
				errorLevel:   internalAuthErrLevel,
			},
		},
	}
}
