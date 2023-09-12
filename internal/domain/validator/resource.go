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
	if agg.GetName() == "" {
		return errors.New("field 'name' cannot be empty")
	}
	if agg.GetFilename() == "" {
		return errors.New("field 'filename' cannot be empty")
	}
	if agg.GetFilepath() == "" {
		return errors.New("field 'path' cannot be empty")
	}
	if agg.GetFilesize() == 0 {
		return errors.New("field 'size' cannot be zero")
	}
	if len(agg.GetFileMIME()) == 0 {
		return errors.New("field 'FileMIME' cannot be empty map")
	}
	return nil
}
