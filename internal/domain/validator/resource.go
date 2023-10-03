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

var supportsFileContentTypes = map[string]struct{}{
	"video/mp4": {},
}

type Resource interface {
	ValidateUploadRequestDTO(req dto.UploadRequest) error
	ValidateEntity(entity entity.Resource) error
	ValidateAggregate(agg *agg.Resource) error
}

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
	if req.GetContentLength() == 0 {
		return errors.NewInvalidUploadedFileError(
			fmt.Sprintf("request file form is empty"),
		)
	}

	contentType := req.GetPart().Header.Get("Content-Type")
	if _, ok := supportsFileContentTypes[contentType]; !ok {
		return errors.NewInvalidUploadedFileError(
			fmt.Sprintf("file '%v' has an unsupported content-type '%v'", req.GetPart().FileName(), contentType),
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
