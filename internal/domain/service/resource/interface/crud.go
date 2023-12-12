package resource_interface

import (
	"github.com/Borislavv/video-streaming/internal/domain/agg"
	dto_interface "github.com/Borislavv/video-streaming/internal/domain/dto/interface"
)

type CRUD interface {
	Upload(reqDTO dto_interface.UploadResourceRequest) (*agg.Resource, error)
	Delete(reqDTO dto_interface.DeleteResourceRequest) (err error)
}
