package builder

import (
	"github.com/Borislavv/video-streaming/internal/domain/agg"
	"github.com/Borislavv/video-streaming/internal/domain/dto"
	"net/http"
)

type Resource interface {
	BuildUploadRequestDTOFromRequest(r *http.Request) (*dto.ResourceUploadRequestDTO, error)
	BuildAggFromUploadRequestDTO(req dto.UploadResourceRequest) *agg.Resource
}
