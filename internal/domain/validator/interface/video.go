package _interface

import (
	"github.com/Borislavv/video-streaming/internal/domain/agg"
	"github.com/Borislavv/video-streaming/internal/domain/dto"
)

type Video interface {
	ValidateGetRequestDTO(req dto.GetVideoRequest) error
	ValidateListRequestDTO(req dto.ListVideoRequest) error
	ValidateCreateRequestDTO(req dto.CreateVideoRequest) error
	ValidateUpdateRequestDTO(req dto.UpdateVideoRequest) error
	ValidateDeleteRequestDTO(req dto.DeleteVideoRequest) error
	ValidateAggregate(agg *agg.Video) error
}
