package builder

import (
	"context"
	"encoding/json"
	"github.com/Borislavv/video-streaming/internal/domain/agg"
	"github.com/Borislavv/video-streaming/internal/domain/dto"
	"github.com/Borislavv/video-streaming/internal/domain/entity"
	"github.com/Borislavv/video-streaming/internal/domain/enum"
	"github.com/Borislavv/video-streaming/internal/domain/errors"
	"github.com/Borislavv/video-streaming/internal/domain/logger"
	"github.com/Borislavv/video-streaming/internal/domain/repository"
	"github.com/Borislavv/video-streaming/internal/domain/service/extractor"
	"github.com/Borislavv/video-streaming/internal/domain/vo"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
	"time"
)

type UserBuilder struct {
	logger         logger.Logger
	ctx            context.Context
	extractor      extractor.RequestParams
	userRepository repository.User
}

// NewUserBuilder is a constructor of UserBuilder
func NewUserBuilder(
	ctx context.Context,
	logger logger.Logger,
	extractor extractor.RequestParams,
	userRepository repository.User,
) *UserBuilder {
	return &UserBuilder{
		ctx:            ctx,
		logger:         logger,
		extractor:      extractor,
		userRepository: userRepository,
	}
}

// BuildGetRequestDTOFromRequest - build a dto.GetUserRequest from raw *http.Request
func (b *UserBuilder) BuildGetRequestDTOFromRequest(r *http.Request) (*dto.UserGetRequestDTO, error) {
	userDTO := &dto.UserGetRequestDTO{}

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

// BuildCreateRequestDTOFromRequest - build a dto.CreateUserRequest from raw *http.Request
func (b *UserBuilder) BuildCreateRequestDTOFromRequest(r *http.Request) (*dto.UserCreateRequestDTO, error) {
	userDTO := &dto.UserCreateRequestDTO{}
	if err := json.NewDecoder(r.Body).Decode(userDTO); err != nil {
		return nil, b.logger.LogPropagate(err)
	}
	return userDTO, nil
}

// BuildAggFromCreateRequestDTO - build an agg.User from dto.CreateUserRequest
func (b *UserBuilder) BuildAggFromCreateRequestDTO(reqDTO dto.CreateUserRequest) (*agg.User, error) {
	// this validation checked previously into the DTO validator
	birthday, err := time.Parse(enum.BirthdayDatePattern, reqDTO.GetBirthday())
	if err != nil {
		// logging the real parsing error
		b.logger.Log(err)
		// logging the error which will be thrown
		return nil, b.logger.LogPropagate(errors.NewBirthdayIsInvalidError(reqDTO.GetBirthday()))
	}

	return &agg.User{
		User: entity.User{
			Username: reqDTO.GetUsername(),
			Password: reqDTO.GetPassword(),
			Email:    reqDTO.GetEmail(),
			Birthday: birthday,
		},
		VideoIDs: []vo.ID{},
		Timestamp: vo.Timestamp{
			CreatedAt: time.Now(),
		},
	}, nil
}

// BuildUpdateRequestDTOFromRequest - build a dto.UserUpdateRequestDTO from raw *http.Request
func (b *UserBuilder) BuildUpdateRequestDTOFromRequest(r *http.Request) (*dto.UserUpdateRequestDTO, error) {
	userDTO := &dto.UserUpdateRequestDTO{}
	if err := json.NewDecoder(r.Body).Decode(&userDTO); err != nil {
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

// BuildAggFromUpdateRequestDTO - build an agg.User from dto.UpdateUserRequest
func (b *UserBuilder) BuildAggFromUpdateRequestDTO(reqDTO dto.UpdateUserRequest) (*agg.User, error) {
	user, err := b.userRepository.Find(b.ctx, reqDTO.GetID())
	if err != nil {
		return nil, b.logger.LogPropagate(err)
	}

	changes := 0
	if reqDTO.GetUsername() != user.Username {
		user.Username = reqDTO.GetUsername()
		changes++
	}

	if reqDTO.GetBirthday() != user.Birthday.String() {
		// this validation checked previously into the DTO validator
		birthday, err := time.Parse(enum.BirthdayDatePattern, reqDTO.GetBirthday())
		if err != nil {
			// here, we must have a valid date or occurred internal error
			return nil, b.logger.CriticalPropagate(err)
		}
		user.Birthday = birthday
		changes++
	}

	if reqDTO.GetPassword() != user.Password {
		user.Password = reqDTO.GetPassword()
		changes++
	}

	if changes > 0 {
		user.Timestamp.UpdatedAt = time.Now()
	}

	return user, nil
}

// BuildDeleteRequestDTOFromRequest - build a dto.DeleteVideoRequest from raw *http.Request
func (b *UserBuilder) BuildDeleteRequestDTOFromRequest(r *http.Request) (*dto.UserDeleteRequestDto, error) {
	getReqDTO, err := b.BuildGetRequestDTOFromRequest(r)
	if err != nil {
		return nil, b.logger.LogPropagate(err)
	}

	return &dto.UserDeleteRequestDto{ID: getReqDTO.ID}, nil
}
