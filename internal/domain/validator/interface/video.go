package validatorinterface

import (
	"github.com/Borislavv/video-streaming/internal/domain/agg"
	"github.com/Borislavv/video-streaming/internal/domain/dto/interface"
)

type Video interface {
	ValidateGetRequestDTO(req dtointerface.GetVideoRequest) error
	ValidateListRequestDTO(req dtointerface.ListVideoRequest) error
	ValidateCreateRequestDTO(req dtointerface.CreateVideoRequest) error
	ValidateUpdateRequestDTO(req dtointerface.UpdateVideoRequest) error
	ValidateDeleteRequestDTO(req dtointerface.DeleteVideoRequest) error
	ValidateAggregate(agg *agg.Video) error
}
