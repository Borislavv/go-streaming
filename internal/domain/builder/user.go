package builder

import (
	"context"
	"encoding/json"
	"github.com/Borislavv/video-streaming/internal/domain/agg"
	"github.com/Borislavv/video-streaming/internal/domain/dto"
	"github.com/Borislavv/video-streaming/internal/domain/entity"
	"github.com/Borislavv/video-streaming/internal/domain/enum"
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
func (b *UserBuilder) BuildAggFromCreateRequestDTO(dto dto.CreateUserRequest) (*agg.User, error) {
	// this validation checked previously into the DTO validator
	birthday, err := time.Parse(enum.BirthdayDatePattern, dto.GetBirthday())
	if err != nil {
		// here, we must have a valid date or occurred internal error
		return nil, b.logger.CriticalPropagate(err)
	}

	return &agg.User{
		User: entity.User{
			Username: dto.GetUsername(),
			Password: dto.GetPassword(),
			Email:    dto.GetEmail(),
			Birthday: birthday,
		},
		VideoIDs: []vo.ID{},
		Timestamp: vo.Timestamp{
			CreatedAt: time.Now(),
		},
	}, nil
}

// BuildAggFromUpdateRequestDTO - build an agg.Video from dto.UpdateVideoRequest
func (b *UserBuilder) BuildAggFromUpdateRequestDTO(dto dto.UpdateUserRequest) (*agg.Video, error) {
	video, err := b.videoRepository.Find(b.ctx, dto.GetID())
	if err != nil {
		return nil, b.logger.LogPropagate(err)
	}

	changes := 0
	if video.Name != dto.GetName() {
		video.Name = dto.GetName()
		changes++
	}
	if video.Description != dto.GetDescription() {
		video.Description = dto.GetDescription()
		changes++
	}
	if !dto.GetResourceID().Value.IsZero() {
		resource, ferr := b.resourceRepository.Find(b.ctx, dto.GetResourceID())
		if ferr != nil {
			return nil, b.logger.LogPropagate(ferr)
		}
		if video.Resource.ID.Value != resource.Resource.ID.Value {
			video.Resource = resource.Resource
			changes++
		}
	}
	if changes > 0 {
		video.Timestamp.UpdatedAt = time.Now()
	}

	return video, nil
}

// BuildDeleteRequestDTOFromRequest - build a dto.DeleteVideoRequest from raw *http.Request
func (b *UserBuilder) BuildDeleteRequestDTOFromRequest(r *http.Request) (*dto.UserDeleteRequestDto, error) {
	getReqDTO, err := b.BuildGetRequestDTOFromRequest(r)
	if err != nil {
		return nil, b.logger.LogPropagate(err)
	}

	return &dto.UserDeleteRequestDto{ID: getReqDTO.ID}, nil
}
