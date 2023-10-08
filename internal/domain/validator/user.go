package validator

import (
	"context"
	"github.com/Borislavv/video-streaming/internal/domain/agg"
	"github.com/Borislavv/video-streaming/internal/domain/dto"
	"github.com/Borislavv/video-streaming/internal/domain/errors"
	"github.com/Borislavv/video-streaming/internal/domain/logger"
	"github.com/Borislavv/video-streaming/internal/domain/repository"
	"github.com/Borislavv/video-streaming/internal/infrastructure/helper"
)

const (
	passwordFiled = "password"
	emailField    = "email"
	birthdayField = "birthday"
)

type UserValidator struct {
	ctx               context.Context
	logger            logger.Logger
	userRepository    repository.User
	adminContactEmail string
}

func NewUserValidator(
	ctx context.Context,
	logger logger.Logger,
	userRepository repository.User,
	adminContactEmail string,
) *UserValidator {
	return &UserValidator{
		ctx:               ctx,
		logger:            logger,
		userRepository:    userRepository,
		adminContactEmail: adminContactEmail,
	}
}

func (v *UserValidator) ValidateGetRequestDTO(reqDTO dto.GetUserRequest) error {
	if reqDTO.GetID().Value.IsZero() {
		return errors.NewFieldCannotBeEmptyError(idField)
	}
	return nil
}

func (v *UserValidator) ValidateCreateRequestDTO(reqDTO dto.CreateUserRequest) error {
	if reqDTO.GetUsername() == "" {
		return errors.NewFieldCannotBeEmptyError(nameField)
	}
	if reqDTO.GetPassword() == "" {
		return errors.NewFieldCannotBeEmptyError(passwordFiled)
	}
	if reqDTO.GetEmail() == "" {
		return errors.NewFieldCannotBeEmptyError(emailField)
	}
	if reqDTO.GetBirthday() == "" {
		return errors.NewFieldCannotBeEmptyError(birthdayField)
	}
	return nil
}

func (v *UserValidator) ValidateUpdateRequestDTO(reqDTO dto.UpdateUserRequest) error {
	if err := v.ValidateGetRequestDTO(reqDTO); err != nil {
		return err
	}
	return nil
}

func (v *UserValidator) ValidateAggregate(agg *agg.User) error {
	// the username cannot be empty or omitted
	if agg.Username == "" {
		return errors.NewFieldCannotBeEmptyError(nameField)
		// the username must be longer than 3 chars and contains only latin letters
	} else if len(agg.Username) < 3 || !helper.IsLatinOnly(agg.Username) {
		return errors.NewUsernameIsInvalidError(agg.Username)
	}

	// the user password cannot be empty or omitted
	if agg.Password == "" {
		return errors.NewFieldCannotBeEmptyError(passwordFiled)
		// the user password must be longer than 8 chars and contains only latin letters/digits
	} else if len(agg.Password) < 8 || !helper.IsLatinOrDigitOnly(agg.Password) {
		return errors.NewPasswordIsInvalidError(agg.Password)
	}

	// the user email cannot be empty or omitted
	if agg.Email == "" {
		return errors.NewFieldCannotBeEmptyError(emailField)
	} else if !helper.IsValidEmail(agg.Email) {
		// logging an email errors for have possibility debug it later
		// when/if a user will report about wrong regex behavior
		err := errors.NewEmailIsInvalidError(agg.Email, v.adminContactEmail)
		return v.logger.WarningPropagate(err)
	}

	// the user birthday cannot be empty or omitted
	if agg.Birthday.IsZero() {
		return errors.NewFieldCannotBeEmptyError(birthdayField)
	}

	// check the user uniqueness
	user, err := v.userRepository.FindByEmail(v.ctx, agg.Email)
	if err != nil {
		if !errors.IsEntityNotFoundError(err) {
			return v.logger.LogPropagate(err)
		}
	}
	if user != nil {
		if !agg.ID.Value.IsZero() {
			if agg.ID.Value.Hex() != user.ID.Value.Hex() {
				//  update user case (check that found user is not the same)
				return errors.NewUserWithSuchEmailAlreadyExistsError(agg.Email)
			}
		} else {
			// create new user case (error thrown if the user was found)
			return errors.NewUserWithSuchEmailAlreadyExistsError(agg.Email)
		}
	}

	return nil
}

func (v *UserValidator) ValidateDeleteRequestDTO(reqDTO dto.DeleteUserRequest) error {
	return v.ValidateGetRequestDTO(reqDTO)
}