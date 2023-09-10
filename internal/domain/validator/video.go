package validator

import (
	"context"
	"errors"
	"github.com/Borislavv/video-streaming/internal/domain/agg"
	"github.com/Borislavv/video-streaming/internal/domain/dto"
	"github.com/Borislavv/video-streaming/internal/domain/errs"
	"github.com/Borislavv/video-streaming/internal/domain/repository"
)

const (
	id   = "id"
	name = "name"
	path = "path"
)

type VideoValidator struct {
	ctx        context.Context
	repository repository.Video
}

func NewVideoValidator(ctx context.Context, repository repository.Video) *VideoValidator {
	return &VideoValidator{ctx: ctx, repository: repository}
}

func (v *VideoValidator) ValidateGetRequestDto(req dto.GetRequest) error {
	if req.GetId().Value.IsZero() {
		return errs.NewFieldCannotBeEmptyError(id)
	}
	return nil
}

func (v *VideoValidator) ValidateListRequestDto(req dto.ListRequest) error {
	if req.GetName() == "" && req.GetPath() == "" {
		return errs.NewAtLeastOneFieldMustBeDefinedError(name, path)
	}
	if req.GetName() != "" && len(req.GetName()) <= 3 {
		return errs.NewFieldLengthMustBeMoreOrLessError(name, true, 3)
	}
	if req.GetPath() != "" && len(req.GetPath()) <= 3 {
		return errs.NewFieldLengthMustBeMoreOrLessError(path, true, 3)
	}
	return nil
}

func (v *VideoValidator) ValidateCreateRequestDto(req dto.CreateRequest) error {
	if req.GetName() == "" {
		return errs.NewFieldCannotBeEmptyError(name)
	}
	if req.GetPath() == "" {
		return errs.NewFieldCannotBeEmptyError(path)
	}
	return nil
}

func (v *VideoValidator) ValidateUpdateRequestDto(req dto.UpdateRequest) error {
	return v.ValidateGetRequestDto(req)
}

func (v *VideoValidator) ValidateDeleteRequestDto(req dto.DeleteRequest) error {
	return v.ValidateGetRequestDto(req)
}

func (v *VideoValidator) ValidateAgg(agg *agg.Video) error {
	if agg.Video.Name == "" {
		return errors.New("'name' cannot be empty")
	}
	if agg.Video.Path == "" {
		return errors.New("'path' cannot be empty")
	}
	has, err := v.repository.Has(v.ctx, agg)
	if err != nil {
		return err
	}
	if has {
		return errs.NewUniquenessCheckFailedError(name, path)
	}
	return nil
}
