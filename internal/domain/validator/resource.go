package validator

import (
	"context"
	"errors"
	"github.com/Borislavv/video-streaming/internal/domain/agg"
	"github.com/Borislavv/video-streaming/internal/domain/dto"
	"github.com/Borislavv/video-streaming/internal/domain/entity"
	"github.com/Borislavv/video-streaming/internal/domain/errs"
	"github.com/Borislavv/video-streaming/internal/domain/repository"
)

type ResourceValidator struct {
	ctx        context.Context
	repository repository.Resource
}

func NewResourceValidator(ctx context.Context, repository repository.Resource) *ResourceValidator {
	return &ResourceValidator{
		ctx:        ctx,
		repository: repository,
	}
}

func (v *ResourceValidator) ValidateUploadRequestDTO(req dto.UploadRequest) error {
	if req.GetHeader().Size == 0 {
		return errs.NewInvalidUploadedFileError(req.GetHeader().Filename)
	}
	return nil
}

func (v *ResourceValidator) ValidateEntity(entity entity.Resource) error {
	if entity.GetName() == "" {
		return errors.New("field 'name' cannot be empty")
	}
	if entity.GetFilename() == "" {
		return errors.New("field 'filename' cannot be empty")
	}
	if entity.GetFilepath() == "" {
		return errors.New("field 'filepath' cannot be empty")
	}
	if entity.GetFilesize() == 0 {
		return errors.New("field 'filesize' cannot be zero")
	}
	if len(entity.GetFileMIME()) == 0 {
		return errors.New("field 'fileMIME' cannot be empty map")
	}
	return nil
}

func (v *ResourceValidator) ValidateAggregate(agg *agg.Resource) error {
	return v.ValidateEntity(agg.Resource)
}
