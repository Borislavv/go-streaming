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

// BuildCreateRequestDtoFromRequest - build a dto.VideoCreateRequestDto from raw *http.Request
func (b *VideoBuilder) BuildCreateRequestDtoFromRequest(r *http.Request) (*dto.VideoCreateRequestDto, error) {
	v := &dto.VideoCreateRequestDto{}
	if err := json.NewDecoder(r.Body).Decode(&v); err != nil {
		return nil, b.logger.LogPropagate(err)
	}
	return v, nil
}

// BuildAggFromCreateRequestDto - build an agg.Video from dto.CreateRequest
func (b *VideoBuilder) BuildAggFromCreateRequestDto(dto dto.CreateRequest) (*agg.Video, error) {
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

// BuildUpdateRequestDtoFromRequest - build a dto.VideoUpdateRequestDto from raw *http.Request
func (b *VideoBuilder) BuildUpdateRequestDtoFromRequest(r *http.Request) (*dto.VideoUpdateRequestDto, error) {
	videoDto := &dto.VideoUpdateRequestDto{}
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

// BuildAggFromUpdateRequestDto - build an agg.Video from dto.UpdateRequest
func (b *VideoBuilder) BuildAggFromUpdateRequestDto(dto dto.UpdateRequest) (*agg.Video, error) {
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
		resource, err := b.resourceRepository.Find(b.ctx, dto.GetResourceID())
		if err != nil {
			return nil, b.logger.LogPropagate(err)
		}
		video.Resource = resource.Resource
		changes++
	}
	if changes > 0 {
		video.Timestamp.UpdatedAt = time.Now()
	}

	return video, nil
}

// BuildGetRequestDtoFromRequest - build a dto.GetRequest from raw *http.Request
func (b *VideoBuilder) BuildGetRequestDtoFromRequest(r *http.Request) (*dto.VideoGetRequestDto, error) {
	videoDto := &dto.VideoGetRequestDto{}

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

// BuildListRequestDtoFromRequest - build a dto.ListRequest from raw *http.Request
func (b *VideoBuilder) BuildListRequestDtoFromRequest(r *http.Request) (*dto.VideoListRequestDto, error) {
	videoDto := &dto.VideoListRequestDto{}

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

// BuildDeleteRequestDtoFromRequest - build a dto.DeleteRequest from raw *http.Request
func (b *VideoBuilder) BuildDeleteRequestDtoFromRequest(r *http.Request) (*dto.VideoDeleteRequestDto, error) {
	videoGetDto, err := b.BuildGetRequestDtoFromRequest(r)
	if err != nil {
		return nil, b.logger.LogPropagate(err)
	}

	return &dto.VideoDeleteRequestDto{ID: videoGetDto.ID}, nil
}
