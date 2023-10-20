package validator

import (
	"github.com/Borislavv/video-streaming/internal/domain/agg"
	"github.com/Borislavv/video-streaming/internal/domain/dto"
	"github.com/Borislavv/video-streaming/internal/domain/entity"
)

type Resource interface {
	// TODO renamed method according to input param name
	ValidateUploadRequestDTO(req dto.UploadResourceRequest) error
	ValidateEntity(entity entity.Resource) error
	ValidateAggregate(agg *agg.Resource) error
	ValidateDeleteRequestDTO(req dto.DeleteResourceRequest) error
}
