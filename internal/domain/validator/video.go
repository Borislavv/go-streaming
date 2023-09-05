package validator

import (
	"errors"
	"github.com/Borislavv/video-streaming/internal/domain/agg"
	"github.com/Borislavv/video-streaming/internal/domain/dto"
	"github.com/Borislavv/video-streaming/internal/domain/errs"
	"github.com/Borislavv/video-streaming/internal/domain/repository"
)

type VideoValidator struct {
	repository repository.Video
}

func NewVideoValidator(repository repository.Video) *VideoValidator {
	return &VideoValidator{repository: repository}
}

func (v *VideoValidator) ValidateCreateRequestDto(dto dto.CreateRequest) error {
	if dto.GetName() == "" {
		return errs.NewFieldCannotBeEmptyError("name")
	}
	if dto.GetPath() == "" {
		return errs.NewFieldCannotBeEmptyError("path")
	}
	return nil
}

func (v *VideoValidator) ValidateUpdateRequestDto(dto dto.UpdateRequest) error {
	//video, found := repository.Find(dto.GetId())
	//if found {
	//	if dto.GetName() == found.Name {
	//		return errs.NewUniquenessCheckFailedError("name")
	//	}
	//}

	if dto.GetId().Value.IsZero() {
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
