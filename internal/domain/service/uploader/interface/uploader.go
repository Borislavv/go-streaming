package uploaderinterface

import (
	"github.com/Borislavv/video-streaming/internal/domain/dto/interface"
)

type Uploader interface {
	// Upload method will be store a file on a disk and calculate a new hashed name. Request DTO mutation!
	Upload(dtointerface.UploadResourceRequest) (err error)
}
