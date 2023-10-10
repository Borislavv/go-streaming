package validator

import (
	"context"
	"github.com/Borislavv/video-streaming/internal/domain/dto"
	"github.com/Borislavv/video-streaming/internal/domain/errors"
	"github.com/Borislavv/video-streaming/internal/domain/logger"
	"github.com/Borislavv/video-streaming/internal/domain/repository"
)

type UserValidator struct {
	ctx                context.Context
	logger             logger.Logger
	resourceValidator  Resource
	videoRepository    repository.Video
	resourceRepository repository.Resource
}

func NewUserValidator(
	ctx context.Context,
	logger logger.Logger,
	resourceValidator Resource,
	videoRepository repository.Video,
	resourceRepository repository.Resource,
) *UserValidator {
	return &UserValidator{
		ctx:                ctx,
		logger:             logger,
		resourceValidator:  resourceValidator,
		videoRepository:    videoRepository,
		resourceRepository: resourceRepository,
	}
}

func (v *UserValidator) ValidateGetRequestDTO(req dto.GetRequest) error {

	if req.GetId().Value.IsZero() {
		return errors.NewFieldCannotBeEmptyError(idField)
	}
	return nil
}
