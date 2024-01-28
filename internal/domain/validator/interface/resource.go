package validatorinterface

import (
	"github.com/Borislavv/video-streaming/internal/domain/agg"
	dtointerface "github.com/Borislavv/video-streaming/internal/domain/dto/interface"
	"github.com/Borislavv/video-streaming/internal/domain/entity"
)

type Resource interface {
	ValidateUploadRequestDTO(req dtointerface.UploadResourceRequest) error
	ValidateEntity(entity entity.Resource) error
	ValidateAggregate(agg *agg.Resource) error
	ValidateDeleteRequestDTO(req dtointerface.DeleteResourceRequest) error
}
