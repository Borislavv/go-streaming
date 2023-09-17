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
	resourceValidator  Resource
	videoRepository    repository.Video
	resourceRepository repository.Resource
}

func NewVideoValidator(
	ctx context.Context,
	logger logger.Logger,
	resourceValidator Resource,
	videoRepository repository.Video,
	resourceRepository repository.Resource,
) *VideoValidator {
	return &VideoValidator{
		ctx:                ctx,
		logger:             logger,
		resourceValidator:  resourceValidator,
		videoRepository:    videoRepository,
		resourceRepository: resourceRepository,
	}
}

func (v *VideoValidator) ValidateGetRequestDTO(req dto.GetRequest) error {
	if req.GetId().Value.IsZero() {
		return errs.NewFieldCannotBeEmptyError(idField)
	}
	return nil
}

func (v *VideoValidator) ValidateListRequestDTO(req dto.ListRequest) error {
	if req.GetName() != "" && len(req.GetName()) <= 3 {
		return errs.NewFieldLengthMustBeMoreOrLessError(nameField, true, 3)
	}
	if !req.GetCreatedAt().IsZero() && (!req.GetFrom().IsZero() || !req.GetTo().IsZero()) {
		return errs.NewValidationError("field 'from' or 'to' cannot be passed with 'createdAt'")
	}
	return nil
}

func (v *VideoValidator) ValidateCreateRequestDTO(req dto.CreateRequest) error {
	if req.GetName() == "" {
		return errs.NewFieldCannotBeEmptyError(nameField)
	}
	if req.GetResourceID().Value.IsZero() {
		return errs.NewFieldCannotBeEmptyError(resourceIDField)
	}
	return nil
}

func (v *VideoValidator) ValidateUpdateRequestDTO(req dto.UpdateRequest) error {
	if err := v.ValidateGetRequestDTO(req); err != nil {
		return err
	}
	return nil
}

func (v *VideoValidator) ValidateDeleteRequestDTO(req dto.DeleteRequest) error {
	return v.ValidateGetRequestDTO(req)
}

func (v *VideoValidator) ValidateAggregate(agg *agg.Video) error {
	// video fields validation
	if agg.Name == "" {
		return errors.New("'name' cannot be empty")
	} else if agg.Resource.ID.Value.IsZero() {
		return errors.New("'resource.id' cannot be empty")
	}

	// resource fields validation
	if err := v.resourceValidator.ValidateEntity(agg.Resource); err != nil {
		return err
	}

	// video validation by name which must be unique
	video, err := v.videoRepository.FindByName(v.ctx, agg.Name)
	if err != nil {
		if !errs.IsNotFoundError(err) {
			return v.logger.LogPropagate(err)
		}
	} else {
		if !agg.ID.Value.IsZero() {
			if video.ID.Value != agg.ID.Value {
				return errs.NewUniquenessCheckFailedError(nameField)
			}
		} else {
			return errs.NewUniquenessCheckFailedError(nameField)
		}
	}

	// video validation by resource.id which must be unique too
	video, err = v.videoRepository.FindByResourceId(v.ctx, agg.Resource.ID)
	if err != nil {
		if !errs.IsNotFoundError(err) {
			return v.logger.LogPropagate(err)
		}
	} else {
		if !agg.ID.Value.IsZero() {
			if video.ID.Value != agg.ID.Value {
				return errs.NewUniquenessCheckFailedError(resourceIDField)
			}
		} else {
			return errs.NewUniquenessCheckFailedError(resourceIDField)
		}
	}

	return nil
}
