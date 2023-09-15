package builder

import (
	"github.com/Borislavv/video-streaming/internal/domain/agg"
	"github.com/Borislavv/video-streaming/internal/domain/dto"
	"github.com/Borislavv/video-streaming/internal/domain/entity"
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

// BuildUploadRequestDtoFromRequest will be parse raw *http.Request and build a dto.UploadRequest
func (b *ResourceBuilder) BuildUploadRequestDtoFromRequest(r *http.Request) (*dto.ResourceUploadRequestDto, error) {
	// request will be parsed and stored in the memory if it is under the RAM threshold,
	// otherwise last parts of parsed file will be stored in the tmp files on the disk space
	if err := r.ParseMultipartForm(b.inMemoryFileSizeThreshold); err != nil {
		return nil, b.logger.LogPropagate(err)
	}

	// receiving a file and header from multipart/form-data
	// by requested filename which is stored in the `formFilename` const.
	file, header, err := r.FormFile(b.formFilename)
	if err != nil {
		return nil, b.logger.LogPropagate(err)
	}
	defer func() { _ = file.Close() }()

	return dto.NewResourceUploadRequest(file, header), nil
}

// BuildAggFromUploadRequestDto will be make an agg.Resource from dto.UploadRequest
func (b *ResourceBuilder) BuildAggFromUploadRequestDto(req dto.UploadRequest) *agg.Resource {
	return &agg.Resource{
		Resource: entity.Resource{
			Name:     req.GetHeader().Filename,
			Filename: req.GetUploadedFilename(),
			Filepath: req.GetUploadedFilepath(),
			Filesize: req.GetHeader().Size,
			FileMIME: req.GetHeader().Header,
		},
		Timestamp: vo.Timestamp{
			CreatedAt: time.Now(),
		},
	}
}
