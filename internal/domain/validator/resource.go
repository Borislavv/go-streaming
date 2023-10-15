package validator

import (
	"context"
	"fmt"
	"github.com/Borislavv/video-streaming/internal/domain/agg"
	"github.com/Borislavv/video-streaming/internal/domain/dto"
	"github.com/Borislavv/video-streaming/internal/domain/entity"
	"github.com/Borislavv/video-streaming/internal/domain/errors"
	"github.com/Borislavv/video-streaming/internal/domain/repository"
)

type ResourceValidator struct {
	ctx         context.Context
	repository  repository.Resource
	maxFilesize int64
}

func NewResourceValidator(ctx context.Context, repository repository.Resource, maxFilesize int64) *ResourceValidator {
	return &ResourceValidator{
		ctx:         ctx,
		repository:  repository,
		maxFilesize: maxFilesize,
	}
}

func (v *ResourceValidator) ValidateUploadRequestDTO(req dto.UploadResourceRequest) error {
	if req.GetRequest().ContentLength == 0 {
		return errors.NewInvalidUploadedFileError("request form file is empty")
	}
	if req.GetRequest().ContentLength > v.maxFilesize {
		return errors.NewInvalidUploadedFileError(
			fmt.Sprintf("request form file is largest than threshold value %d", v.maxFilesize),
		)
	}
	return nil
}

func (v *ResourceValidator) ValidateEntity(entity entity.Resource) error {
	if entity.GetName() == "" {
		return errors.NewInternalValidationError("field 'name' cannot be empty")
	}
	if entity.GetFilename() == "" {
		return errors.NewInternalValidationError("field 'filename' cannot be empty")
	}
	if entity.GetFilepath() == "" {
		return errors.NewInternalValidationError("field 'filepath' cannot be empty")
	}
	if entity.GetFilesize() == 0 {
		return errors.NewInternalValidationError("field 'filesize' cannot be zero")
	}
	if entity.GetFiletype() == "" {
		return errors.NewInternalValidationError("field 'filetype' cannot be empty")
	}
	return nil
}

func (v *ResourceValidator) ValidateAggregate(agg *agg.Resource) error {
	return v.ValidateEntity(agg.Resource)
}

func (v *ResourceValidator) ValidateGetRequestDTO(req dto.GetResourceRequest) error {
	if req.GetId().Value.IsZero() {
		return errors.NewFieldCannotBeEmptyError(idField)
	}
	return nil
}

func (v *ResourceValidator) ValidateDeleteRequestDTO(req dto.DeleteResourceRequest) error {
	return v.ValidateGetRequestDTO(req)
}
