package validator

import (
	"context"
	"fmt"
	"github.com/Borislavv/video-streaming/internal/domain/agg"
	dtointerface "github.com/Borislavv/video-streaming/internal/domain/dto/interface"
	"github.com/Borislavv/video-streaming/internal/domain/entity"
	"github.com/Borislavv/video-streaming/internal/domain/errtype"
	repositoryinterface "github.com/Borislavv/video-streaming/internal/domain/repository/interface"
	diinterface "github.com/Borislavv/video-streaming/internal/domain/service/di/interface"
)

type ResourceValidator struct {
	ctx         context.Context
	repository  repositoryinterface.Resource
	maxFilesize int64
}

func NewResourceValidator(serviceContainer diinterface.ServiceContainer) (*ResourceValidator, error) {
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

func (v *ResourceValidator) ValidateUploadRequestDTO(req dtointerface.UploadResourceRequest) error {
	if req.GetUserID().Value.IsZero() {
		return errtype.NewFieldCannotBeEmptyError(userIDField)
	}
	if req.GetRequest().ContentLength == 0 {
		return errtype.NewInvalidUploadedFileError("request form file is empty")
	}
	if req.GetRequest().ContentLength > v.maxFilesize {
		return errtype.NewInvalidUploadedFileError(
			fmt.Sprintf("request form file is largest than threshold value %d", v.maxFilesize),
		)
	}
	return nil
}

func (v *ResourceValidator) ValidateEntity(entity entity.Resource) error {
	if entity.UserID.Value.IsZero() {
		return errtype.NewInternalValidationError("field 'userID' cannot be empty")
	}
	if entity.GetName() == "" {
		return errtype.NewInternalValidationError("field 'name' cannot be empty")
	}
	if entity.GetFilename() == "" {
		return errtype.NewInternalValidationError("field 'filename' cannot be empty")
	}
	if entity.GetFilepath() == "" {
		return errtype.NewInternalValidationError("field 'filepath' cannot be empty")
	}
	if entity.GetFilesize() == 0 {
		return errtype.NewInternalValidationError("field 'filesize' cannot be zero")
	}
	if entity.GetFiletype() == "" {
		return errtype.NewInternalValidationError("field 'filetype' cannot be empty")
	}
	return nil
}

func (v *ResourceValidator) ValidateAggregate(agg *agg.Resource) error {
	return v.ValidateEntity(agg.Resource)
}

func (v *ResourceValidator) ValidateGetRequestDTO(req dtointerface.GetResourceRequest) error {
	if req.GetID().Value.IsZero() {
		return errtype.NewFieldCannotBeEmptyError(idField)
	}
	if req.GetUserID().Value.IsZero() {
		return errtype.NewFieldCannotBeEmptyError(userIDField)
	}
	return nil
}

func (v *ResourceValidator) ValidateDeleteRequestDTO(req dtointerface.DeleteResourceRequest) error {
	return v.ValidateGetRequestDTO(req)
}
