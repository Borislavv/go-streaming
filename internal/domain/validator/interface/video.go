package validator_interface

import (
	"github.com/Borislavv/video-streaming/internal/domain/agg"
	"github.com/Borislavv/video-streaming/internal/domain/dto/interface"
)

type Video interface {
	ValidateGetRequestDTO(req dto_interface.GetVideoRequest) error
	ValidateListRequestDTO(req dto_interface.ListVideoRequest) error
	ValidateCreateRequestDTO(req dto_interface.CreateVideoRequest) error
	ValidateUpdateRequestDTO(req dto_interface.UpdateVideoRequest) error
	ValidateDeleteRequestDTO(req dto_interface.DeleteVideoRequest) error
	ValidateAggregate(agg *agg.Video) error
}
