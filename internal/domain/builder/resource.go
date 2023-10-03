package builder

import (
	"github.com/Borislavv/video-streaming/internal/domain/agg"
	"github.com/Borislavv/video-streaming/internal/domain/dto"
	"github.com/Borislavv/video-streaming/internal/domain/entity"
	"github.com/Borislavv/video-streaming/internal/domain/errors"
	"github.com/Borislavv/video-streaming/internal/domain/logger"
	"github.com/Borislavv/video-streaming/internal/domain/vo"
	"io"
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

// BuildUploadRequestDTOFromRequest will be parse raw *http.Request and build a dto.UploadRequest
func (b *ResourceBuilder) BuildUploadRequestDTOFromRequest(r *http.Request) (*dto.ResourceUploadRequestDTO, error) {
	// extract the multipart form reader (handling the form as a stream)
	reader, err := r.MultipartReader()
	if err != nil {
		return nil, b.logger.LogPropagate(err)
	}

	for { // search the file part
		part, err := reader.NextPart()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, b.logger.ErrorPropagate(err)
		}

		if part.FileName() != "" {
			return dto.NewResourceUploadRequest(part, r.ContentLength), nil
		}
	}

	return nil, errors.NewFormDoesNotContainsUploadedFileError()
}

// BuildAggFromUploadRequestDTO will be make an agg.Resource from dto.UploadRequest
func (b *ResourceBuilder) BuildAggFromUploadRequestDTO(req dto.UploadRequest) *agg.Resource {
	return &agg.Resource{
		Resource: entity.Resource{
			Name:     req.GetPart().FileName(),
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
