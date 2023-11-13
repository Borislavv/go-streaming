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

func (s *CRUDService) Get(req dto.GetUserRequest) (user *agg.User, err error) {
	if err = s.validator.ValidateGetRequestDTO(req); err != nil {
		return nil, s.logger.LogPropagate(err)
	}

	if !req.GetID().Value.IsZero() {
		user, err = s.repository.FindOneByID(s.ctx, req)
		if err != nil {
			return nil, s.logger.LogPropagate(err)
		}
	} else if req.GetEmail() != "" {
		user, err = s.repository.FindOneByEmail(s.ctx, req)
		if err != nil {
			return nil, s.logger.LogPropagate(err)
		}
	}

	return user, nil
}

func (s *CRUDService) Create(req dto.CreateUserRequest) (*agg.User, error) {
	// validation of input request
	if err := s.validator.ValidateCreateRequestDTO(req); err != nil {
		return nil, s.logger.LogPropagate(err)
	}

	// building an aggregate
	userAgg, err := s.builder.BuildAggFromCreateRequestDTO(req)
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

func (s *CRUDService) Update(req dto.UpdateUserRequest) (*agg.User, error) {
	// validation of input request
	if err := s.validator.ValidateUpdateRequestDTO(req); err != nil {
		return nil, s.logger.LogPropagate(err)
	}

	// building an aggregate
	userAgg, err := s.builder.BuildAggFromUpdateRequestDTO(req)
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

func (s *CRUDService) Delete(req dto.DeleteUserRequest) (err error) {
	// validation of input request
	if err = s.validator.ValidateDeleteRequestDTO(req); err != nil {
		return s.logger.LogPropagate(err)
	}

	// fetching a user which will be deleted
	userAgg, err := s.repository.FindOneByID(s.ctx, req)
	if err != nil {
		return s.logger.LogPropagate(err)
	}

	// fetching a video list which will be deleted
	videoAggs, total, err := s.videoService.List(&dto.VideoListRequestDTO{UserID: userAgg.ID})
	if err != nil {
		return s.logger.LogPropagate(err)
	}

	// removing the references video first
	if total > 0 {
		for _, videoAgg := range videoAggs {
			if err = s.videoService.Delete(dto.NewVideoDeleteRequestDto(videoAgg.ID, userAgg.ID)); err != nil {
				if errors.IsEntityNotFoundError(err) {
					s.logger.Warning(
						fmt.Sprintf("user delete warning: reference video '%v' is not exists", videoAgg.ID.Value),
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
