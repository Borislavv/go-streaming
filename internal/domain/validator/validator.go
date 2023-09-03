package validator

import (
	"github.com/Borislavv/video-streaming/internal/domain/agg"
	"github.com/Borislavv/video-streaming/internal/domain/dto"
)

type Video interface {
	ValidateCreateRequestDto(dto dto.CreateRequest) error
	ValidateUpdateRequestDto(dto dto.UpdateRequest) error
	ValidateAgg(agg *agg.Video) error
}
