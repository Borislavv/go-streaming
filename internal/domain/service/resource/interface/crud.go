package resourceinterface

import (
	"github.com/Borislavv/video-streaming/internal/domain/agg"
	dtointerface "github.com/Borislavv/video-streaming/internal/domain/dto/interface"
)

type CRUD interface {
	Upload(reqDTO dtointerface.UploadResourceRequest) (*agg.Resource, error)
	Delete(reqDTO dtointerface.DeleteResourceRequest) (err error)
}
