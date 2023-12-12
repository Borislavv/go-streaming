package builder_interface

import (
	"github.com/Borislavv/video-streaming/internal/domain/agg"
	"github.com/Borislavv/video-streaming/internal/domain/dto"
	"github.com/Borislavv/video-streaming/internal/domain/dto/interface"
	"net/http"
)

type Resource interface {
	BuildUploadRequestDTOFromRequest(r *http.Request) (*dto.ResourceUploadRequestDTO, error)
	BuildAggFromUploadRequestDTO(reqDTO dto_interface.UploadResourceRequest) *agg.Resource
}
