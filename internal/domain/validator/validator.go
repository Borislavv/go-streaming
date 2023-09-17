package validator

import (
	"github.com/Borislavv/video-streaming/internal/domain/agg"
	"github.com/Borislavv/video-streaming/internal/domain/dto"
	"github.com/Borislavv/video-streaming/internal/domain/entity"
)

type Resource interface {
	ValidateUploadRequestDTO(req dto.UploadRequest) error
	ValidateEntity(entity entity.Resource) error
	ValidateAggregate(agg *agg.Resource) error
}

type Video interface {
	ValidateGetRequestDTO(req dto.GetRequest) error
	ValidateListRequestDTO(req dto.ListRequest) error
	ValidateCreateRequestDTO(req dto.CreateRequest) error
	ValidateUpdateRequestDTO(req dto.UpdateRequest) error
	ValidateDeleteRequestDTO(req dto.DeleteRequest) error
	ValidateAggregate(agg *agg.Video) error
}
