package validator

import (
	"context"
	"errors"
	"github.com/Borislavv/video-streaming/internal/domain/agg"
	"github.com/Borislavv/video-streaming/internal/domain/dto"
	"github.com/Borislavv/video-streaming/internal/domain/errs"
	"github.com/Borislavv/video-streaming/internal/domain/repository"
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
		return errs.NewFieldCannotBeEmptyError("id")
	}
	return nil
}

func (v *VideoValidator) ValidateCreateRequestDto(req dto.CreateRequest) error {
	if req.GetName() == "" {
		return errs.NewFieldCannotBeEmptyError("name")
	}
	if req.GetPath() == "" {
		return errs.NewFieldCannotBeEmptyError("path")
	}
	return nil
}

func (v *VideoValidator) ValidateUpdateRequestDto(req dto.UpdateRequest) error {
	if req.GetId().Value.IsZero() {
		return errs.NewFieldCannotBeEmptyError("id")
	}
	return nil
}

func (v *VideoValidator) ValidateAgg(agg *agg.Video) error {
	if agg.Video.Name == "" {
		return errors.New("'name' cannot be empty")
	}
	if agg.Video.Path == "" {
		return errors.New("'path' cannot be empty")
	}
	return nil
}
