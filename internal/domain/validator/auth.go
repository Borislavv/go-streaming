package validator

import (
	"github.com/Borislavv/video-streaming/internal/domain/dto"
	"github.com/Borislavv/video-streaming/internal/domain/enum"
	"github.com/Borislavv/video-streaming/internal/domain/errors"
	"github.com/Borislavv/video-streaming/internal/domain/logger"
	"net/http"
)

type AuthValidator struct {
	logger                   logger.Logger
	adminContactEmailAddress string
}

func NewAuthValidator(
	logger logger.Logger,
	adminContactEmailAddress string,
) *AuthValidator {
	return &AuthValidator{
		logger:                   logger,
		adminContactEmailAddress: adminContactEmailAddress,
	}
}

// ValidateAuthRequest is method which will check the auth request DTO on valid.
func (v *AuthValidator) ValidateAuthRequest(req dto.AuthRequest) error {
	if req.GetEmail() == "" {
		return errors.NewFieldCannotBeEmptyError(emailField)
	}

	if req.GetPassword() == "" {
		return errors.NewFieldCannotBeEmptyError(passwordField)
	}

	return nil
}

// ValidateTokennessRequest is method which will check that access token header exists.
func (v *AuthValidator) ValidateTokennessRequest(r *http.Request) error {
	if token := r.Header.Get(enum.AccessTokenHeaderKey); token != "" {
		return nil
	}

	if _, err := r.Cookie(enum.AccessTokenHeaderKey); err == nil {
		return nil
	}

	return errors.NewAuthFailedError("token is not provided")
}
