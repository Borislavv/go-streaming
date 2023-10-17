package resource

import (
	"github.com/Borislavv/video-streaming/internal/domain/agg"
	"github.com/Borislavv/video-streaming/internal/domain/dto"
)

type CRUD interface {
	Upload(reqDTO dto.UploadResourceRequest) (*agg.Resource, error)
	Delete(reqDTO dto.DeleteResourceRequest) (err error)
}
