package validator_interface

import (
	"github.com/Borislavv/video-streaming/internal/domain/agg"
	dto_interface "github.com/Borislavv/video-streaming/internal/domain/dto/interface"
	"github.com/Borislavv/video-streaming/internal/domain/entity"
)

type Resource interface {
	ValidateUploadRequestDTO(req dto_interface.UploadResourceRequest) error
	ValidateEntity(entity entity.Resource) error
	ValidateAggregate(agg *agg.Resource) error
	ValidateDeleteRequestDTO(req dto_interface.DeleteResourceRequest) error
}
