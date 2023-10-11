package user

import (
	"context"
	"github.com/Borislavv/video-streaming/internal/domain/agg"
	"github.com/Borislavv/video-streaming/internal/domain/builder"
	"github.com/Borislavv/video-streaming/internal/domain/dto"
	"github.com/Borislavv/video-streaming/internal/domain/logger"
	"github.com/Borislavv/video-streaming/internal/domain/repository"
	"github.com/Borislavv/video-streaming/internal/domain/validator"
)

type CRUDService struct {
	ctx        context.Context
	logger     logger.Logger
	builder    builder.User
	validator  validator.User
	repository repository.User
}

func NewCRUDService(
	ctx context.Context,
	logger logger.Logger,
	builder builder.User,
	validator validator.User,
	repository repository.User,
) *CRUDService {
	return &CRUDService{
		ctx:        ctx,
		logger:     logger,
		builder:    builder,
		validator:  validator,
		repository: repository,
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
