package user

import (
	"context"
	"fmt"
	"github.com/Borislavv/video-streaming/internal/domain/agg"
	"github.com/Borislavv/video-streaming/internal/domain/builder/interface"
	"github.com/Borislavv/video-streaming/internal/domain/dto"
	dto_interface "github.com/Borislavv/video-streaming/internal/domain/dto/interface"
	"github.com/Borislavv/video-streaming/internal/domain/errors"
	"github.com/Borislavv/video-streaming/internal/domain/logger/interface"
	repository_interface "github.com/Borislavv/video-streaming/internal/domain/repository/interface"
	di_interface "github.com/Borislavv/video-streaming/internal/domain/service/di/interface"
	video_interface "github.com/Borislavv/video-streaming/internal/domain/service/video/interface"
	validator_interface "github.com/Borislavv/video-streaming/internal/domain/validator/interface"
)

type CRUDService struct {
	ctx          context.Context
	logger       logger_interface.Logger
	builder      builder_interface.User
	validator    validator_interface.User
	repository   repository_interface.User
	videoService video_interface.CRUD
}

func NewCRUDService(serviceContainer di_interface.ContainerManager) (*CRUDService, error) {
	loggerService, err := serviceContainer.GetLoggerService()
	if err != nil {
		return nil, err
	}

	ctx, err := serviceContainer.GetCtx()
	if err != nil {
		return nil, loggerService.LogPropagate(err)
	}

	userBuilder, err := serviceContainer.GetUserBuilder()
	if err != nil {
		return nil, loggerService.LogPropagate(err)
	}

	userValidator, err := serviceContainer.GetUserValidator()
	if err != nil {
		return nil, loggerService.LogPropagate(err)
	}

	userRepository, err := serviceContainer.GetUserRepository()
	if err != nil {
		return nil, loggerService.LogPropagate(err)
	}

	videoCRUDService, err := serviceContainer.GetVideoCRUDService()
	if err != nil {
		return nil, loggerService.LogPropagate(err)
	}

	return &CRUDService{
		ctx:          ctx,
		logger:       loggerService,
		builder:      userBuilder,
		validator:    userValidator,
		repository:   userRepository,
		videoService: videoCRUDService,
	}, nil
}

func (s *CRUDService) Get(req dto_interface.GetUserRequest) (user *agg.User, err error) {
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

func (s *CRUDService) Create(req dto_interface.CreateUserRequest) (*agg.User, error) {
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

func (s *CRUDService) Update(req dto_interface.UpdateUserRequest) (*agg.User, error) {
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

func (s *CRUDService) Delete(req dto_interface.DeleteUserRequest) (err error) {
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
