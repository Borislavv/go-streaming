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

func (s *CRUDService) Get(req dto.GetRequest) (user *agg.User, err error) {
	if err = s.validator.ValidateGetRequestDTO(req); err != nil {
		return nil, s.logger.LogPropagate(err)
	}

	user, err = s.repository.Find(s.ctx, req.GetId())
	if err != nil {
		return nil, s.logger.LogPropagate(err)
	}

	return user, nil
}
