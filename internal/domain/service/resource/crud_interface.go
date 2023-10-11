package resource

import (
	"github.com/Borislavv/video-streaming/internal/domain/agg"
	"github.com/Borislavv/video-streaming/internal/domain/dto"
)

type CRUD interface {
	Upload(req dto.UploadResourceRequest) (*agg.Resource, error)
}
