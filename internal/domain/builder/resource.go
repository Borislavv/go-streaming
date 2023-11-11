package builder

import (
	"github.com/Borislavv/video-streaming/internal/domain/agg"
	"github.com/Borislavv/video-streaming/internal/domain/dto"
	"github.com/Borislavv/video-streaming/internal/domain/entity"
	"github.com/Borislavv/video-streaming/internal/domain/enum"
	"github.com/Borislavv/video-streaming/internal/domain/logger"
	"github.com/Borislavv/video-streaming/internal/domain/vo"
	"net/http"
	"time"
)

type ResourceBuilder struct {
	logger                    logger.Logger
	formFilename              string
	inMemoryFileSizeThreshold int64
}

func NewResourceBuilder(
	logger logger.Logger,
	formFilename string,
	inMemoryFileSizeThreshold int64,
) *ResourceBuilder {
	return &ResourceBuilder{
		logger:                    logger,
		formFilename:              formFilename,
		inMemoryFileSizeThreshold: inMemoryFileSizeThreshold,
	}
}

// BuildUploadRequestDTOFromRequest will be parse raw *http.Request and build a dto.UploadResourceRequest
func (b *ResourceBuilder) BuildUploadRequestDTOFromRequest(r *http.Request) (*dto.ResourceUploadRequestDTO, error) {
	resourceDTO := dto.NewResourceUploadRequest(r)
	if userID, ok := r.Context().Value(enum.UserIDContextKey).(vo.ID); ok {
		resourceDTO.SetUserID(userID)
	}
	return resourceDTO, nil
}

// BuildAggFromUploadRequestDTO will be make an agg.Resource from dto.UploadResourceRequest
func (b *ResourceBuilder) BuildAggFromUploadRequestDTO(req dto.UploadResourceRequest) *agg.Resource {
	return &agg.Resource{
		Resource: entity.Resource{
			UserID:   req.GetUserID(),
			Name:     req.GetOriginFilename(),
			Filename: req.GetUploadedFilename(),
			Filepath: req.GetUploadedFilepath(),
			Filesize: req.GetUploadedFilesize(),
			Filetype: req.GetUploadedFiletype(),
		},
		Timestamp: vo.Timestamp{
			CreatedAt: time.Now(),
		},
	}
}
