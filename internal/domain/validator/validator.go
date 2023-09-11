package validator

import (
	"github.com/Borislavv/video-streaming/internal/domain/agg"
	"github.com/Borislavv/video-streaming/internal/domain/dto"
)

type Resource interface {
	ValidateUploadRequestDto(req dto.UploadRequest) error
	ValidateAgg(agg *agg.Resource) error
}

type Video interface {
	ValidateGetRequestDto(req dto.GetRequest) error
	ValidateListRequestDto(req dto.ListRequest) error
	ValidateCreateRequestDto(req dto.CreateRequest) error
	ValidateUpdateRequestDto(req dto.UpdateRequest) error
	ValidateDeleteRequestDto(req dto.DeleteRequest) error
	ValidateAgg(agg *agg.Video) error
}
