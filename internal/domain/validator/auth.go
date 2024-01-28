package validator

import (
	"github.com/Borislavv/video-streaming/internal/domain/dto/interface"
	"github.com/Borislavv/video-streaming/internal/domain/enum"
	"github.com/Borislavv/video-streaming/internal/domain/errors"
	"github.com/Borislavv/video-streaming/internal/domain/logger/interface"
	diinterface "github.com/Borislavv/video-streaming/internal/domain/service/di/interface"
	"net/http"
)

type AuthValidator struct {
	logger                   loggerinterface.Logger
	adminContactEmailAddress string
}

func NewAuthValidator(serviceContainer diinterface.ContainerManager) (*AuthValidator, error) {
	loggerService, err := serviceContainer.GetLoggerService()
	if err != nil {
		return nil, err
	}

	cfg, err := serviceContainer.GetConfig()
	if err != nil {
		return nil, loggerService.LogPropagate(err)
	}

	return &AuthValidator{
		logger:                   loggerService,
		adminContactEmailAddress: cfg.AdminContactEmail,
	}, nil
}

// ValidateAuthRequest is method which will check the auth request DTO on valid.
func (v *AuthValidator) ValidateAuthRequest(req dtointerface.AuthRequest) error {
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
