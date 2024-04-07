package builder

import (
	"context"
	"encoding/json"
	"github.com/Borislavv/video-streaming/internal/domain/agg"
	"github.com/Borislavv/video-streaming/internal/domain/dto"
	dtointerface "github.com/Borislavv/video-streaming/internal/domain/dto/interface"
	"github.com/Borislavv/video-streaming/internal/domain/entity"
	"github.com/Borislavv/video-streaming/internal/domain/enum"
	"github.com/Borislavv/video-streaming/internal/domain/errtype"
	"github.com/Borislavv/video-streaming/internal/domain/logger/interface"
	repositoryinterface "github.com/Borislavv/video-streaming/internal/domain/repository/interface"
	"github.com/Borislavv/video-streaming/internal/domain/service/di/interface"
	extractorinterface "github.com/Borislavv/video-streaming/internal/domain/service/extractor/interface"
	securityinterface "github.com/Borislavv/video-streaming/internal/domain/service/security/interface"
	"github.com/Borislavv/video-streaming/internal/domain/vo"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"io"
	"net/http"
	"time"
)

type UserBuilder struct {
	logger         loggerinterface.Logger
	ctx            context.Context
	extractor      extractorinterface.RequestParams
	userRepository repositoryinterface.User
	passwordHasher securityinterface.PasswordHasher
}

// NewUserBuilder is a constructor of UserBuilder.
func NewUserBuilder(serviceContainer diinterface.ContainerManager) (*UserBuilder, error) {
	loggerService, err := serviceContainer.GetLoggerService()
	if err != nil {
		return nil, err
	}

	ctx, err := serviceContainer.GetCtx()
	if err != nil {
		return nil, loggerService.LogPropagate(err)
	}

	requestParametersExtractorService, err := serviceContainer.GetRequestParametersExtractorService()
	if err != nil {
		return nil, loggerService.LogPropagate(err)
	}

	userRepository, err := serviceContainer.GetUserRepository()
	if err != nil {
		return nil, loggerService.LogPropagate(err)
	}

	passwordHasherService, err := serviceContainer.GetPasswordHasherService()
	if err != nil {
		return nil, loggerService.LogPropagate(err)
	}

	return &UserBuilder{
		ctx:            ctx,
		logger:         loggerService,
		extractor:      requestParametersExtractorService,
		userRepository: userRepository,
		passwordHasher: passwordHasherService,
	}, nil
}

// BuildGetRequestDTOFromRequest - build a dto.GetUserRequest from raw *http.Request.
func (b *UserBuilder) BuildGetRequestDTOFromRequest(r *http.Request) (*dto.UserGetRequestDTO, error) {
	userDTO := &dto.UserGetRequestDTO{}
	if userID, ok := r.Context().Value(enum.UserIDContextKey).(vo.ID); ok {
		userDTO.ID = userID
	}
	return userDTO, nil
}

// BuildCreateRequestDTOFromRequest - build a dto.CreateUserRequest from raw *http.Request.
func (b *UserBuilder) BuildCreateRequestDTOFromRequest(r *http.Request) (*dto.UserCreateRequestDTO, error) {
	userDTO := &dto.UserCreateRequestDTO{}
	if err := json.NewDecoder(r.Body).Decode(userDTO); err != nil {
		if err == io.EOF {
			return nil, b.logger.LogPropagate(errtype.NewRequestBodyIsEmptyError())
		}
		return nil, b.logger.LogPropagate(err)
	}
	return userDTO, nil
}

// BuildAggFromCreateRequestDTO - build an agg.User from dto.CreateUserRequest
func (b *UserBuilder) BuildAggFromCreateRequestDTO(req dtointerface.CreateUserRequest) (*agg.User, error) {
	// this validation checked previously into the DTO validator
	birthday, err := time.Parse(enum.BirthdayDatePattern, req.GetBirthday())
	if err != nil {
		// logging the real parsing error
		b.logger.Log(err)
		// logging the error which will be thrown
		return nil, b.logger.LogPropagate(errtype.NewBirthdayIsInvalidError(req.GetBirthday()))
	}

	// hash user's real password
	passwordHash, err := b.passwordHasher.Hash(req.GetPassword())
	if err != nil {
		return nil, b.logger.LogPropagate(err)
	}

	u := &agg.User{
		User: entity.User{
			Username: req.GetUsername(),
			Email:    req.GetEmail(),
			Birthday: birthday,
		},
		Timestamp: vo.Timestamp{
			CreatedAt: time.Now(),
		},
	}
	u.SetPassword(passwordHash)
	return u, nil
}

// BuildUpdateRequestDTOFromRequest - build a dto.UserUpdateRequestDTO from raw *http.Request.
func (b *UserBuilder) BuildUpdateRequestDTOFromRequest(r *http.Request) (*dto.UserUpdateRequestDTO, error) {
	userDTO := &dto.UserUpdateRequestDTO{}
	if err := json.NewDecoder(r.Body).Decode(&userDTO); err != nil {
		if err == io.EOF {
			return nil, b.logger.LogPropagate(errtype.NewRequestBodyIsEmptyError())
		}
		return nil, b.logger.LogPropagate(err)
	}

	hexID, err := b.extractor.GetParameter(idField, r)
	if err != nil {
		return nil, b.logger.LogPropagate(err)
	}
	oID, err := primitive.ObjectIDFromHex(hexID)
	if err != nil {
		return nil, b.logger.LogPropagate(err)
	}
	userDTO.ID = vo.ID{Value: oID}

	return userDTO, nil
}

// BuildAggFromUpdateRequestDTO - build an agg.User from dto.UpdateUserRequest.
func (b *UserBuilder) BuildAggFromUpdateRequestDTO(req dtointerface.UpdateUserRequest) (*agg.User, error) {
	user, err := b.userRepository.FindOneByID(b.ctx, req)
	if err != nil {
		return nil, b.logger.LogPropagate(err)
	}

	changes := 0
	if req.GetUsername() != user.Username {
		user.Username = req.GetUsername()
		changes++
	}

	if req.GetBirthday() != user.Birthday.String() {
		// this validation checked previously into the DTO validator
		birthday, perr := time.Parse(enum.BirthdayDatePattern, req.GetBirthday())
		if perr != nil {
			// here, we must have a valid date or occurred internal error
			return nil, b.logger.CriticalPropagate(perr)
		}
		user.Birthday = birthday
		changes++
	}

	// hash user's real password
	passwordHash, err := b.passwordHasher.Hash(req.GetPassword())
	if err != nil {
		return nil, b.logger.LogPropagate(err)
	}

	if req.GetPassword() != passwordHash {
		user.SetPassword(passwordHash)
		changes++
	}

	if changes > 0 {
		user.Timestamp.UpdatedAt = time.Now()
	}

	return user, nil
}

// BuildDeleteRequestDTOFromRequest - build a dto.DeleteVideoRequest from raw *http.Request.
func (b *UserBuilder) BuildDeleteRequestDTOFromRequest(r *http.Request) (*dto.UserDeleteRequestDTO, error) {
	getReqDTO, err := b.BuildGetRequestDTOFromRequest(r)
	if err != nil {
		return nil, b.logger.LogPropagate(err)
	}

	return &dto.UserDeleteRequestDTO{ID: getReqDTO.ID}, nil
}

func (b *UserBuilder) BuildResponseDTO(user *agg.User) (*dto.UserResponseDTO, error) {
	return &dto.UserResponseDTO{
		ID:       user.ID,
		Username: user.Username,
		Email:    user.Email,
		Birthday: user.Birthday,
	}, nil
}
