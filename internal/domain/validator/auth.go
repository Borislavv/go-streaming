package validator

import (
	"github.com/Borislavv/video-streaming/internal/domain/dto"
	"github.com/Borislavv/video-streaming/internal/domain/errors"
	"github.com/Borislavv/video-streaming/internal/domain/logger"
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

func (v *AuthValidator) ValidateAuthRequest(reqDTO dto.AuthRequest) error {
	if reqDTO.GetEmail() == "" {
		return errors.NewFieldCannotBeEmptyError(emailField)
	}
	if reqDTO.GetPassword() == "" {
		return errors.NewFieldCannotBeEmptyError(passwordFiled)
	}
	return nil
}
