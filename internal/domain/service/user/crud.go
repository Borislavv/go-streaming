package user

import (
	"context"
	"fmt"
	"github.com/Borislavv/video-streaming/internal/domain/agg"
	"github.com/Borislavv/video-streaming/internal/domain/builder"
	"github.com/Borislavv/video-streaming/internal/domain/dto"
	"github.com/Borislavv/video-streaming/internal/domain/errors"
	"github.com/Borislavv/video-streaming/internal/domain/logger"
	"github.com/Borislavv/video-streaming/internal/domain/repository"
	"github.com/Borislavv/video-streaming/internal/domain/service/video"
	"github.com/Borislavv/video-streaming/internal/domain/validator"
)

type CRUDService struct {
	ctx          context.Context
	logger       logger.Logger
	builder      builder.User
	validator    validator.User
	repository   repository.User
	videoService video.CRUD
}

func NewCRUDService(
	ctx context.Context,
	logger logger.Logger,
	builder builder.User,
	validator validator.User,
	repository repository.User,
	videoService video.CRUD,
) *CRUDService {
	return &CRUDService{
		ctx:          ctx,
		logger:       logger,
		builder:      builder,
		validator:    validator,
		repository:   repository,
		videoService: videoService,
	}
}

func (s *CRUDService) Get(reqDTO dto.GetUserRequest) (user *agg.User, err error) {
	if err = s.validator.ValidateGetRequestDTO(reqDTO); err != nil {
		return nil, s.logger.LogPropagate(err)
	}

	user, err = s.repository.Find(s.ctx, reqDTO.GetID())
	if err != nil {
		return nil, s.logger.LogPropagate(err)
	}

	return user, nil
}

func (s *CRUDService) Create(reqDTO dto.CreateUserRequest) (*agg.User, error) {
	// validation of input request
	if err := s.validator.ValidateCreateRequestDTO(reqDTO); err != nil {
		return nil, s.logger.LogPropagate(err)
	}

	// building an aggregate
	userAgg, err := s.builder.BuildAggFromCreateRequestDTO(reqDTO)
	if err != nil {
		return nil, s.logger.LogPropagate(err)
	}

	// validation of an aggregate
	if err = s.validator.ValidateAggregate(userAgg); err != nil {
		return nil, s.logger.LogPropagate(err)
	}

	// saving an aggregate into storage
	userAgg, err = s.repository.Insert(s.ctx, userAgg)
	if err != nil {
		return nil, s.logger.LogPropagate(err)
	}

	return userAgg, nil
}

func (s *CRUDService) Update(reqDTO dto.UpdateUserRequest) (*agg.User, error) {
	// validation of input request
	if err := s.validator.ValidateUpdateRequestDTO(reqDTO); err != nil {
		return nil, s.logger.LogPropagate(err)
	}

	// building an aggregate
	userAgg, err := s.builder.BuildAggFromUpdateRequestDTO(reqDTO)
	if err != nil {
		return nil, s.logger.LogPropagate(err)
	}

	// validation an aggregate
	if err = s.validator.ValidateAggregate(userAgg); err != nil {
		return nil, s.logger.LogPropagate(err)
	}

	// saving the updated aggregate into storage
	userAgg, err = s.repository.Update(s.ctx, userAgg)
	if err != nil {
		return nil, s.logger.LogPropagate(err)
	}

	return userAgg, nil
}

func (s *CRUDService) Delete(reqDTO dto.DeleteUserRequest) (err error) {
	// validation of input request
	if err = s.validator.ValidateDeleteRequestDTO(reqDTO); err != nil {
		return s.logger.LogPropagate(err)
	}

	// fetching a user which will be deleted
	userAgg, err := s.repository.Find(s.ctx, reqDTO.GetID())
	if err != nil {
		return s.logger.LogPropagate(err)
	}

	// removing the references video first
	if len(userAgg.VideoIDs) > 0 {
		for _, videoID := range userAgg.VideoIDs {
			if err = s.videoService.Delete(&dto.VideoDeleteRequestDto{ID: videoID}); err != nil {
				if errors.IsEntityNotFoundError(err) {
					s.logger.Warning(
						fmt.Sprintf("user remobing error: references video '%v' is not exists", videoID.Value),
					)
					continue
				}
			}
		}
	}

	// user removing
	if err = s.repository.Remove(s.ctx, userAgg); err != nil {
		return s.logger.LogPropagate(err)
	}

	return nil
}
