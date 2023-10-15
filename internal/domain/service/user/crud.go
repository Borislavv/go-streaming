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
	"github.com/Borislavv/video-streaming/internal/domain/validator"
)

type CRUDService struct {
	ctx             context.Context
	logger          logger.Logger
	builder         builder.User
	validator       validator.User
	repository      repository.User
	videoRepository repository.Video
}

func NewCRUDService(
	ctx context.Context,
	logger logger.Logger,
	builder builder.User,
	validator validator.User,
	repository repository.User,
	videoRepository repository.Video,
) *CRUDService {
	return &CRUDService{
		ctx:             ctx,
		logger:          logger,
		builder:         builder,
		validator:       validator,
		repository:      repository,
		videoRepository: videoRepository,
	}
}

func (s *CRUDService) Get(req dto.GetUserRequest) (user *agg.User, err error) {
	if err = s.validator.ValidateGetRequestDTO(req); err != nil {
		return nil, s.logger.LogPropagate(err)
	}

	user, err = s.repository.Find(s.ctx, req.GetId())
	if err != nil {
		return nil, s.logger.LogPropagate(err)
	}

	return user, nil
}

func (s *CRUDService) Create(userDTO dto.CreateUserRequest) (*agg.User, error) {
	// validation of input request
	if err := s.validator.ValidateCreateRequestDTO(userDTO); err != nil {
		return nil, s.logger.LogPropagate(err)
	}

	// building an aggregate
	userAgg, err := s.builder.BuildAggFromCreateRequestDTO(userDTO)
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

func (s *CRUDService) Delete(reqDTO dto.DeleteUserRequest) error {
	// validation of input request
	if err := s.validator.ValidateDeleteRequestDTO(reqDTO); err != nil {
		return s.logger.LogPropagate(err)
	}

	// fetching a user which will be deleted
	userAgg, err := s.repository.Find(s.ctx, reqDTO.GetId())
	if err != nil {
		return s.logger.LogPropagate(err)
	}

	// removing the references video first
	if len(userAgg.VideoIDs) > 0 {
		for _, videoID := range userAgg.VideoIDs {
			// fetching the removing video
			videoAgg, err := s.videoRepository.Find(s.ctx, videoID)
			if err != nil {
				if errors.IsEntityNotFoundError(err) {
					s.logger.Warning(
						fmt.Sprintf("user remobing error: references video '%v' is not exists", videoID.Value),
					)
					continue
				}
				return s.logger.LogPropagate(err)
			}
			// removing the fetched video
			if err := s.videoRepository.Remove(s.ctx, videoAgg); err != nil {
				return s.logger.LogPropagate(err)
			}
		}
	}

	// user removing
	if err = s.repository.Remove(s.ctx, userAgg); err != nil {
		return s.logger.LogPropagate(err)
	}

	return nil
}
