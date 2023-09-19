package builder

import (
	"context"
	"encoding/json"
	"github.com/Borislavv/video-streaming/internal/domain/agg"
	"github.com/Borislavv/video-streaming/internal/domain/api/request"
	"github.com/Borislavv/video-streaming/internal/domain/dto"
	"github.com/Borislavv/video-streaming/internal/domain/entity"
	"github.com/Borislavv/video-streaming/internal/domain/logger"
	"github.com/Borislavv/video-streaming/internal/domain/repository"
	"github.com/Borislavv/video-streaming/internal/domain/vo"
	"github.com/Borislavv/video-streaming/internal/infrastructure/helper"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
	"strconv"
	"time"
)

const (
	idField           = "id"
	nameField         = "name"
	createdAtField    = "createdAt"
	fromField         = "from"
	toField           = "to"
	pageField         = "page"
	limitField        = "limit"
	limitDefaultValue = 25
	pageDefaultValue  = 1
)

type VideoBuilder struct {
	logger             logger.Logger
	ctx                context.Context
	extractor          request.Extractor
	videoRepository    repository.Video
	resourceRepository repository.Resource
}

// NewVideoBuilder is a constructor of VideoBuilder
func NewVideoBuilder(
	ctx context.Context,
	logger logger.Logger,
	extractor request.Extractor,
	videoRepository repository.Video,
	resourceRepository repository.Resource,
) *VideoBuilder {
	return &VideoBuilder{
		ctx:                ctx,
		logger:             logger,
		extractor:          extractor,
		videoRepository:    videoRepository,
		resourceRepository: resourceRepository,
	}
}

// BuildCreateRequestDTOFromRequest - build a dto.CreateRequest from raw *http.Request
func (b *VideoBuilder) BuildCreateRequestDTOFromRequest(r *http.Request) (*dto.VideoCreateRequestDTO, error) {
	v := &dto.VideoCreateRequestDTO{}
	if err := json.NewDecoder(r.Body).Decode(&v); err != nil {
		return nil, b.logger.LogPropagate(err)
	}
	return v, nil
}

// BuildAggFromCreateRequestDTO - build an agg.Video from dto.CreateRequest
func (b *VideoBuilder) BuildAggFromCreateRequestDTO(dto dto.CreateRequest) (*agg.Video, error) {
	resource, err := b.resourceRepository.Find(b.ctx, dto.GetResourceID())
	if err != nil {
		return nil, b.logger.LogPropagate(err)
	}

	return &agg.Video{
		Video: entity.Video{
			Name:        dto.GetName(),
			Description: dto.GetDescription(),
		},
		Resource: resource.Resource,
		Timestamp: vo.Timestamp{
			CreatedAt: time.Now(),
		},
	}, nil
}

// BuildUpdateRequestDTOFromRequest - build a dto.UpdateRequest from raw *http.Request
func (b *VideoBuilder) BuildUpdateRequestDTOFromRequest(r *http.Request) (*dto.VideoUpdateRequestDTO, error) {
	videoDto := &dto.VideoUpdateRequestDTO{}
	if err := json.NewDecoder(r.Body).Decode(&videoDto); err != nil {
		return nil, b.logger.LogPropagate(err)
	}

	hexId, err := b.extractor.GetParameter(idField, r)
	if err != nil {
		return nil, b.logger.LogPropagate(err)
	}
	oid, err := primitive.ObjectIDFromHex(hexId)
	if err != nil {
		return nil, b.logger.LogPropagate(err)
	}
	videoDto.ID = vo.ID{Value: oid}

	return videoDto, nil
}

// BuildAggFromUpdateRequestDTO - build an agg.Video from dto.UpdateRequest
func (b *VideoBuilder) BuildAggFromUpdateRequestDTO(dto dto.UpdateRequest) (*agg.Video, error) {
	video, err := b.videoRepository.Find(b.ctx, dto.GetId())
	if err != nil {
		return nil, b.logger.LogPropagate(err)
	}

	changes := 0
	if video.Name != dto.GetName() {
		video.Name = dto.GetName()
		changes++
	}
	if video.Description != dto.GetDescription() {
		video.Description = dto.GetDescription()
		changes++
	}
	if !dto.GetResourceID().Value.IsZero() {
		resource, ferr := b.resourceRepository.Find(b.ctx, dto.GetResourceID())
		if ferr != nil {
			return nil, b.logger.LogPropagate(ferr)
		}
		if video.Resource.ID.Value != resource.Resource.ID.Value {
			video.Resource = resource.Resource
			changes++
		}
	}
	if changes > 0 {
		video.Timestamp.UpdatedAt = time.Now()
	}

	return video, nil
}

// BuildGetRequestDTOFromRequest - build a dto.GetRequest from raw *http.Request
func (b *VideoBuilder) BuildGetRequestDTOFromRequest(r *http.Request) (*dto.VideoGetRequestDTO, error) {
	videoDto := &dto.VideoGetRequestDTO{}

	hexId, err := b.extractor.GetParameter(idField, r)
	if err != nil {
		return nil, b.logger.LogPropagate(err)
	}
	oid, err := primitive.ObjectIDFromHex(hexId)
	if err != nil {
		return nil, b.logger.LogPropagate(err)
	}
	videoDto.ID = vo.ID{Value: oid}

	return videoDto, nil
}

// BuildListRequestDTOFromRequest - build a dto.ListRequest from raw *http.Request
func (b *VideoBuilder) BuildListRequestDTOFromRequest(r *http.Request) (*dto.VideoListRequestDTO, error) {
	videoDto := &dto.VideoListRequestDTO{}

	if b.extractor.HasParameter(nameField, r) {
		if nm, err := b.extractor.GetParameter(nameField, r); err == nil {
			videoDto.Name = nm
		}
	}
	if b.extractor.HasParameter(createdAtField, r) {
		createdAt, _ := b.extractor.GetParameter(createdAtField, r)

		parsedCreatedAt, err := helper.ParseTime(createdAt)
		if err != nil {
			return nil, b.logger.LogPropagate(err)
		} else {
			videoDto.CreatedAt = parsedCreatedAt
		}
	}
	if b.extractor.HasParameter(fromField, r) {
		from, _ := b.extractor.GetParameter(fromField, r)

		parsedFrom, err := helper.ParseTime(from)
		if err != nil {
			return nil, b.logger.LogPropagate(err)
		} else {
			videoDto.From = parsedFrom
		}
	}
	if b.extractor.HasParameter(toField, r) {
		to, _ := b.extractor.GetParameter(toField, r)

		parsedTo, err := helper.ParseTime(to)
		if err != nil {
			return nil, b.logger.LogPropagate(err)
		} else {
			videoDto.To = parsedTo
		}
	}
	if b.extractor.HasParameter(pageField, r) {
		pg, _ := b.extractor.GetParameter(pageField, r)
		pgi, atoiErr := strconv.Atoi(pg)
		if atoiErr != nil {
			return nil, b.logger.LogPropagate(atoiErr)
		}
		videoDto.Page = pgi
	} else {
		videoDto.Page = pageDefaultValue
	}
	if b.extractor.HasParameter(limitField, r) {
		l, _ := b.extractor.GetParameter(limitField, r)
		li, atoiErr := strconv.Atoi(l)
		if atoiErr != nil {
			return nil, b.logger.LogPropagate(atoiErr)
		}
		videoDto.Limit = li
	} else {
		videoDto.Limit = limitDefaultValue
	}

	return videoDto, nil
}

// BuildDeleteRequestDTOFromRequest - build a dto.DeleteRequest from raw *http.Request
func (b *VideoBuilder) BuildDeleteRequestDTOFromRequest(r *http.Request) (*dto.VideoDeleteRequestDto, error) {
	videoGetDto, err := b.BuildGetRequestDTOFromRequest(r)
	if err != nil {
		return nil, b.logger.LogPropagate(err)
	}

	return &dto.VideoDeleteRequestDto{ID: videoGetDto.ID}, nil
}
