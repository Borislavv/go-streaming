package validator

import (
	"context"
	"errors"
	"github.com/Borislavv/video-streaming/internal/domain/agg"
	"github.com/Borislavv/video-streaming/internal/domain/dto"
	"github.com/Borislavv/video-streaming/internal/domain/errs"
	"github.com/Borislavv/video-streaming/internal/domain/logger"
	"github.com/Borislavv/video-streaming/internal/domain/repository"
)

const (
	idField         = "id"
	nameField       = "name"
	resourceIDField = "resourceID"
)

type VideoValidator struct {
	ctx                context.Context
	logger             logger.Logger
	videoRepository    repository.Video
	resourceRepository repository.Resource
}

func NewVideoValidator(
	ctx context.Context,
	logger logger.Logger,
	videoRepository repository.Video,
	resourceRepository repository.Resource,
) *VideoValidator {
	return &VideoValidator{
		ctx:                ctx,
		logger:             logger,
		videoRepository:    videoRepository,
		resourceRepository: resourceRepository,
	}
}

func (v *VideoValidator) ValidateGetRequestDto(req dto.GetRequest) error {
	if req.GetId().Value.IsZero() {
		return errs.NewFieldCannotBeEmptyError(idField)
	}
	return nil
}

func (v *VideoValidator) ValidateListRequestDto(req dto.ListRequest) error {
	if req.GetName() != "" && len(req.GetName()) <= 3 {
		return errs.NewFieldLengthMustBeMoreOrLessError(nameField, true, 3)
	}
	if !req.GetCreatedAt().IsZero() && (!req.GetFrom().IsZero() || !req.GetTo().IsZero()) {
		return errs.NewValidationError("field 'from' or 'to' cannot be passed with 'createdAt'")
	}
	return nil
}

func (v *VideoValidator) ValidateCreateRequestDto(req dto.CreateRequest) error {
	if req.GetName() == "" {
		return errs.NewFieldCannotBeEmptyError(nameField)
	}
	if req.GetResourceID().Value.IsZero() {
		return errs.NewFieldCannotBeEmptyError(resourceIDField)
	}
	return nil
}

func (v *VideoValidator) ValidateUpdateRequestDto(req dto.UpdateRequest) error {
	if err := v.ValidateGetRequestDto(req); err != nil {
		return err
	}

	if !req.GetResourceID().Value.IsZero() {
		resource, err := v.resourceRepository.Find(v.ctx, req.GetResourceID())
		if err != nil {
			return err
		}

		video, err := v.videoRepository.FindByResource(v.ctx, resource)
		if err != nil {
			if errs.IsNotFoundError(err) {
				return nil
			}
			return v.logger.LogPropagate(err)
		}

		if video.ID == req.GetId() {
			return errs.NewUniquenessCheckFailedError(resourceIDField)
		}
	}

	return nil
}

func (v *VideoValidator) ValidateDeleteRequestDto(req dto.DeleteRequest) error {
	return v.ValidateGetRequestDto(req)
}

func (v *VideoValidator) ValidateAgg(agg *agg.Video) error {
	if agg.Name == "" {
		return errors.New("'name' cannot be empty")
	}

	has, err := v.videoRepository.Has(v.ctx, agg)
	if err != nil {
		return v.logger.LogPropagate(err)
	}
	if has {
		return errs.NewUniquenessCheckFailedError(nameField, resourceIDField)
	}

	return nil
}
