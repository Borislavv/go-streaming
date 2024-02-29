package builderinterface

import (
	"github.com/Borislavv/video-streaming/internal/domain/agg"
	"github.com/Borislavv/video-streaming/internal/domain/dto"
	dtointerface "github.com/Borislavv/video-streaming/internal/domain/dto/interface"
	"net/http"
)

type Resource interface {
	BuildUploadRequestDTOFromRequest(r *http.Request) (*dto.ResourceUploadRequestDTO, error)
	BuildAggFromUploadRequestDTO(reqDTO dtointerface.UploadResourceRequest) *agg.Resource
}
