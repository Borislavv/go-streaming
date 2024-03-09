package builder

import (
	"github.com/Borislavv/video-streaming/internal/domain/agg"
	"github.com/Borislavv/video-streaming/internal/domain/dto"
	"github.com/Borislavv/video-streaming/internal/domain/dto/interface"
	"github.com/Borislavv/video-streaming/internal/domain/entity"
	"github.com/Borislavv/video-streaming/internal/domain/enum"
	"github.com/Borislavv/video-streaming/internal/domain/logger/interface"
	"github.com/Borislavv/video-streaming/internal/domain/service/di/interface"
	"github.com/Borislavv/video-streaming/internal/domain/vo"
	"net/http"
	"time"
)

type ResourceBuilder struct {
	logger                    loggerinterface.Logger
	formFilename              string
	inMemoryFileSizeThreshold int64
}

func NewResourceBuilder(serviceContainer diinterface.ContainerManager) (*ResourceBuilder, error) {
	loggerService, err := serviceContainer.GetLoggerService()
	if err != nil {
		return nil, err
	}

	cfg, err := serviceContainer.GetConfig()
	if err != nil {
		return nil, loggerService.LogPropagate(err)
	}

	return &ResourceBuilder{
		logger:                    loggerService,
		formFilename:              cfg.ResourceFormFilename,
		inMemoryFileSizeThreshold: cfg.ResourceInMemoryFileSizeThreshold,
	}, nil
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
func (b *ResourceBuilder) BuildAggFromUploadRequestDTO(req dto_interface.UploadResourceRequest) *agg.Resource {
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
