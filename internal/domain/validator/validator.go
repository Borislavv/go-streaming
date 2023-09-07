package validator

import (
	"github.com/Borislavv/video-streaming/internal/domain/agg"
	"github.com/Borislavv/video-streaming/internal/domain/dto"
)

type Video interface {
	ValidateGetRequestDto(req dto.GetRequest) error
	ValidateCreateRequestDto(req dto.CreateRequest) error
	ValidateUpdateRequestDto(req dto.UpdateRequest) error
	ValidateDeleteRequestDto(req dto.DeleteRequest) error
	ValidateAgg(agg *agg.Video) error
}
