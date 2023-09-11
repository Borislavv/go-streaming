package validator

import (
	"context"
	"errors"
	"github.com/Borislavv/video-streaming/internal/domain/agg"
	"github.com/Borislavv/video-streaming/internal/domain/dto"
	"github.com/Borislavv/video-streaming/internal/domain/errs"
	"github.com/Borislavv/video-streaming/internal/domain/repository"
)

type ResourceValidator struct {
	ctx        context.Context
	repository repository.Resource
}

func NewResourceValidator(
	ctx context.Context,
	repository repository.Resource,
) *ResourceValidator {
	return &ResourceValidator{
		ctx:        ctx,
		repository: repository,
	}
}

func (v *ResourceValidator) ValidateUploadRequestDto(req dto.UploadRequest) error {
	if req.GetHeader().Size == 0 {
		return errs.NewInvalidUploadedFileError(req.GetHeader().Filename)
	}
	return nil
}

func (v *ResourceValidator) ValidateAgg(agg *agg.Resource) error {
	if agg.Resource.GetName() == "" {
		return errors.New("field 'name' cannot be empty")
	}
	if agg.Resource.GetFilename() == "" {
		return errors.New("field 'filename' cannot be empty")
	}
	if agg.Resource.GetPath() == "" {
		return errors.New("field 'path' cannot be empty")
	}
	if agg.Resource.GetSize() == 0 {
		return errors.New("field 'size' cannot be zero")
	}
	if len(agg.Resource.GetMIME()) == 0 {
		return errors.New("field 'MIME' cannot be empty map")
	}
	return nil
}
