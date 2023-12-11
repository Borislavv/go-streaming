package validator

import (
	"context"
	"fmt"
	"github.com/Borislavv/video-streaming/internal/domain/agg"
	"github.com/Borislavv/video-streaming/internal/domain/dto/interface"
	"github.com/Borislavv/video-streaming/internal/domain/entity"
	"github.com/Borislavv/video-streaming/internal/domain/errors"
	repository_interface "github.com/Borislavv/video-streaming/internal/domain/repository/interface"
	di_interface "github.com/Borislavv/video-streaming/internal/domain/service/di/interface"
)

type ResourceValidator struct {
	ctx         context.Context
	repository  repository_interface.Resource
	maxFilesize int64
}

func NewResourceValidator(serviceContainer di_interface.ContainerManager) (*ResourceValidator, error) {
	loggerService, err := serviceContainer.GetLoggerService()
	if err != nil {
		return nil, err
	}

	ctx, err := serviceContainer.GetCtx()
	if err != nil {
		return nil, loggerService.LogPropagate(err)
	}

	repo, err := serviceContainer.GetResourceRepository()
	if err != nil {
		return nil, loggerService.LogPropagate(err)
	}

	cfg, err := serviceContainer.GetConfig()
	if err != nil {
		return nil, loggerService.LogPropagate(err)
	}

	return &ResourceValidator{
		ctx:         ctx,
		repository:  repo,
		maxFilesize: cfg.ResourceMaxFilesizeThreshold,
	}, nil
}

func (v *ResourceValidator) ValidateUploadRequestDTO(req dto_interface.UploadResourceRequest) error {
	if req.GetUserID().Value.IsZero() {
		return errors.NewFieldCannotBeEmptyError(userIDField)
	}
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
	if entity.UserID.Value.IsZero() {
		return errors.NewInternalValidationError("field 'userID' cannot be empty")
	}
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

func (v *ResourceValidator) ValidateGetRequestDTO(req dto_interface.GetResourceRequest) error {
	if req.GetID().Value.IsZero() {
		return errors.NewFieldCannotBeEmptyError(idField)
	}
	if req.GetUserID().Value.IsZero() {
		return errors.NewFieldCannotBeEmptyError(userIDField)
	}
	return nil
}

func (v *ResourceValidator) ValidateDeleteRequestDTO(req dto_interface.DeleteResourceRequest) error {
	return v.ValidateGetRequestDTO(req)
}
