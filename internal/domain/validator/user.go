package validator

import (
	"context"
	"github.com/Borislavv/video-streaming/internal/domain/agg"
	"github.com/Borislavv/video-streaming/internal/domain/dto"
	dto_interface "github.com/Borislavv/video-streaming/internal/domain/dto/interface"
	"github.com/Borislavv/video-streaming/internal/domain/enum"
	"github.com/Borislavv/video-streaming/internal/domain/errors"
	"github.com/Borislavv/video-streaming/internal/domain/logger/interface"
	repository_interface "github.com/Borislavv/video-streaming/internal/domain/repository/interface"
	diinterface "github.com/Borislavv/video-streaming/internal/domain/service/di/interface"
	"github.com/Borislavv/video-streaming/internal/domain/vo"
	"github.com/Borislavv/video-streaming/internal/infrastructure/helper"
	"time"
)

const (
	emailField    = "email"
	passwordField = "password"
	birthdayField = "birthday"
)

type UserValidator struct {
	ctx               context.Context
	logger            logger_interface.Logger
	userRepository    repository_interface.User
	adminContactEmail string
}

func NewUserValidator(serviceContainer diinterface.ContainerManager) (*UserValidator, error) {
	loggerService, err := serviceContainer.GetLoggerService()
	if err != nil {
		return nil, err
	}

	ctx, err := serviceContainer.GetCtx()
	if err != nil {
		return nil, loggerService.LogPropagate(err)
	}

	userRepository, err := serviceContainer.GetUserRepository()
	if err != nil {
		return nil, loggerService.LogPropagate(err)
	}

	cfg, err := serviceContainer.GetConfig()
	if err != nil {
		return nil, loggerService.LogPropagate(err)
	}

	return &UserValidator{
		ctx:               ctx,
		logger:            loggerService,
		userRepository:    userRepository,
		adminContactEmail: cfg.AdminContactEmail,
	}, nil
}

func (v *UserValidator) ValidateGetRequestDTO(req dto_interface.GetUserRequest) error {
	if req.GetID().Value.IsZero() && req.GetEmail() == "" {
		return errors.NewAtLeastOneFieldMustBeDefinedError(idField, emailField)
	}
	return nil
}

func (v *UserValidator) ValidateCreateRequestDTO(req dto_interface.CreateUserRequest) (err error) {
	if err = v.isValidUsername(req.GetUsername()); err != nil {
		return v.logger.LogPropagate(err)
	}

	if err = v.isValidPassword(req.GetPassword()); err != nil {
		return v.logger.LogPropagate(err)
	}

	if err = v.isValidEmail(req.GetEmail()); err != nil {
		return v.logger.LogPropagate(err)
	}

	if err = v.isValidBirthday(req.GetBirthday()); err != nil {
		return v.logger.LogPropagate(err)
	}

	if err = v.isUniqueUser(req.GetEmail()); err != nil {
		return v.logger.LogPropagate(err)
	}

	return nil
}

func (v *UserValidator) ValidateUpdateRequestDTO(req dto_interface.UpdateUserRequest) (err error) {
	if req.GetID().Value.IsZero() {
		return errors.NewFieldCannotBeEmptyError(idField)
	}

	if err = v.isValidUsername(req.GetUsername()); err != nil {
		return v.logger.LogPropagate(err)
	}

	if err = v.isValidPassword(req.GetPassword()); err != nil {
		return v.logger.LogPropagate(err)
	}

	if err = v.isValidBirthday(req.GetBirthday()); err != nil {
		return v.logger.LogPropagate(err)
	}

	return nil
}

func (v *UserValidator) ValidateAggregate(agg *agg.User) error {
	// the username cannot be empty or omitted
	if agg.Username == "" {
		return errors.NewInternalValidationError("user agg was built with empty username")
	}

	// the user password cannot be empty or omitted
	if agg.GetPassword() == "" {
		return errors.NewInternalValidationError("user agg was built with empty password")
	}

	// the user email cannot be empty or omitted
	if agg.Email == "" {
		return errors.NewInternalValidationError("user agg was built with empty email")
	}

	// the user birthday cannot be empty or omitted
	if agg.Birthday.IsZero() {
		return errors.NewInternalValidationError("user agg was built with empty birthday")
	}

	return nil
}

func (v *UserValidator) ValidateDeleteRequestDTO(req dto_interface.DeleteUserRequest) error {
	if req.GetID().Value.IsZero() {
		return errors.NewFieldCannotBeEmptyError(idField)
	}
	return nil
}

func (v *UserValidator) isValidBirthday(birthday string) error {
	if birthday == "" {
		return errors.NewFieldCannotBeEmptyError(birthdayField)
	}

	_, err := time.Parse(enum.BirthdayDatePattern, birthday)
	if err != nil {
		return errors.NewBirthdayIsInvalidError(birthday)
	}

	return nil
}

func (v *UserValidator) isValidEmail(email string) error {
	if email == "" {
		return errors.NewFieldCannotBeEmptyError(emailField)
	}

	// logging an email errors for have possibility debug it later
	// when/if a user will report about wrong regex behavior
	if !helper.IsValidEmail(email) {
		return errors.NewEmailIsInvalidError(email, v.adminContactEmail)
	}

	return nil
}

func (v *UserValidator) isValidPassword(password string) error {
	if password == "" {
		return errors.NewFieldCannotBeEmptyError(passwordField)
	}

	// the user password must be longer than 8 chars and contains only latin letters/digits
	if len(password) < 8 || !helper.IsLatinOrDigitOnly(password) {
		return errors.NewPasswordIsInvalidError(password)
	}

	return nil
}

func (v *UserValidator) isValidUsername(username string) error {
	if username == "" {
		return errors.NewFieldCannotBeEmptyError(nameField)
	}

	// the username must be longer than 3 chars and contains only latin letters
	if len(username) < 3 || !helper.IsLatinOnly(username) {
		return errors.NewUsernameIsInvalidError(username)
	}

	return nil
}

// isUniqueUser checks whether an email is unique per user collection.
func (v *UserValidator) isUniqueUser(email string) error {
	// UserGetRequestDTO must be created with specifying the email only otherwise a user will be found by id in any case
	user, err := v.userRepository.FindOneByEmail(v.ctx, dto.NewUserGetRequestDTO(vo.ID{}, email))
	if err != nil && !errors.IsEntityNotFoundError(err) {
		return v.logger.LogPropagate(err)
	}

	if user != nil {
		return errors.NewUserWithSuchEmailAlreadyExistsError(email)
	}

	return nil
}
