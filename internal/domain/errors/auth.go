package errors

import (
	"fmt"
	"github.com/Borislavv/video-streaming/internal/infrastructure/service/logger"
	"net/http"
)

const (
	authErrType           = "authorization"
	internalAuthErrLevel  = logger.CriticalLevel
	internalAuthErrStatus = http.StatusInternalServerError
	publicAuthErrLevel    = logger.InfoLevel
	publicAuthErrStatus   = http.StatusBadRequest
)

type AuthFailedError struct{ publicError }

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
				errorLevel:   publicAuthErrLevel,
				errorStatus:  publicAuthErrStatus,
			},
		},
	}
}

type AccessTokenIsEmptyOrOmittedError struct{ publicError }

func NewAccessTokenIsEmptyOrOmittedError() *AccessTokenIsEmptyOrOmittedError {
	return &AccessTokenIsEmptyOrOmittedError{
		publicError{
			errored{
				ErrorMessage: "authorization failed: token is empty or omitted",
				ErrorType:    authErrType,
				errorStatus:  publicAuthErrStatus,
				errorLevel:   publicAuthErrLevel,
			},
		},
	}
}

type AccessTokenIsInvalidError struct{ publicError }

func NewAccessTokenIsInvalidError() *AccessTokenIsInvalidError {
	return &AccessTokenIsInvalidError{
		publicError{
			errored{
				ErrorMessage: "authorization failed: provided token is invalid",
				ErrorType:    authErrType,
				errorStatus:  publicAuthErrStatus,
				errorLevel:   publicAuthErrLevel,
			},
		},
	}
}

type AccessTokenWasBlockedError struct{ publicError }

func NewAccessTokenWasBlockedError() *AccessTokenWasBlockedError {
	return &AccessTokenWasBlockedError{
		publicError{
			errored{
				ErrorMessage: "authorization failed: token was blocked",
				ErrorType:    authErrType,
				errorStatus:  publicAuthErrStatus,
				errorLevel:   publicAuthErrLevel,
			},
		},
	}
}

type TokenAlgoWasNotMatchedError struct{ internalError }

func NewTokenAlgoWasNotMatchedInternalError(token string) *TokenAlgoWasNotMatchedError {
	return &TokenAlgoWasNotMatchedError{
		internalError{
			errored{
				ErrorMessage: fmt.Sprintf("token '%v' algo was not matched", token),
				ErrorType:    authErrType,
				errorStatus:  internalAuthErrStatus,
				errorLevel:   internalAuthErrLevel,
			},
		},
	}
}

type TokenUnexpectedSigningMethodError struct{ internalError }

func NewTokenUnexpectedSigningMethodInternalError(token string, algo interface{}) *TokenUnexpectedSigningMethodError {
	return &TokenUnexpectedSigningMethodError{
		internalError{
			errored{
				ErrorMessage: fmt.Sprintf("unexpected signing algo '%v' for token '%v'", algo, token),
				ErrorType:    authErrType,
				errorStatus:  internalAuthErrStatus,
				errorLevel:   internalAuthErrLevel,
			},
		},
	}
}

type TokenInvalidError struct{ internalError }

func NewTokenInvalidInternalError(token string) *TokenInvalidError {
	return &TokenInvalidError{
		internalError{
			errored{
				ErrorMessage: fmt.Sprintf("token '%v' is not valid", token),
				ErrorType:    authErrType,
				errorStatus:  internalAuthErrStatus,
				errorLevel:   internalAuthErrLevel,
			},
		},
	}
}

type TokenIssuerWasNotMatchedError struct{ internalError }

func NewTokenIssuerWasNotMatchedInternalError(token string) *TokenIssuerWasNotMatchedError {
	return &TokenIssuerWasNotMatchedError{
		internalError{
			errored{
				ErrorMessage: fmt.Sprintf("token '%v' issuer was not matched", token),
				ErrorType:    authErrType,
				errorStatus:  internalAuthErrStatus,
				errorLevel:   internalAuthErrLevel,
			},
		},
	}
}
