package validator

import (
	"errors"
	"github.com/Borislavv/video-streaming/internal/domain/agg"
	"github.com/Borislavv/video-streaming/internal/domain/dto"
	"github.com/Borislavv/video-streaming/internal/errs"
)

type VideoValidator struct {
}

func NewVideoValidator() *VideoValidator {
	return &VideoValidator{}
}

func (v *VideoValidator) ValidateCreateRequestDto(dto dto.CreateRequest) error {
	if dto.GetName() == "" {
		return errs.NewFieldCannotBeEmptyError("'name' field cannot be empty or omitted")
	}
	if dto.GetPath() == "" {
		return errs.NewFieldCannotBeEmptyError("'path' field cannot be empty or omitted")
	}
	return nil
}

func (v *VideoValidator) ValidateUpdateRequestDto(dto dto.UpdateRequest) error {
	if dto.GetId().Value.IsZero() {
		return errs.NewFieldCannotBeEmptyError("'id' cannot be empty or omitted")
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
