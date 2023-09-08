package builder

import (
	"context"
	"encoding/json"
	"github.com/Borislavv/video-streaming/internal/domain/agg"
	"github.com/Borislavv/video-streaming/internal/domain/dto"
	"github.com/Borislavv/video-streaming/internal/domain/entity"
	"github.com/Borislavv/video-streaming/internal/domain/repository"
	"github.com/Borislavv/video-streaming/internal/domain/vo"
	"github.com/Borislavv/video-streaming/internal/infrastructure/api/v1/request"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
	"strconv"
	"time"
)

const (
	id                = "id"
	name              = "name"
	path              = "path"
	page              = "page"
	limit             = "limit"
	limitDefaultValue = 25
	pageDefaultValue  = 1
)

type VideoBuilder struct {
	ctx        context.Context
	extractor  request.Extractor // TODO must be removed due to DDD (infrastructure leaked into the domain logic)
	repository repository.Video
}

func NewVideoBuilder(ctx context.Context, extractor request.Extractor, repository repository.Video) *VideoBuilder {
	return &VideoBuilder{ctx: ctx, extractor: extractor, repository: repository}
}

// BuildCreateRequestDtoFromRequest - build a dto.VideoCreateRequestDto from raw *http.Request
func (b *VideoBuilder) BuildCreateRequestDtoFromRequest(r *http.Request) (*dto.VideoCreateRequestDto, error) {
	v := &dto.VideoCreateRequestDto{}
	if err := json.NewDecoder(r.Body).Decode(&v); err != nil {
		return nil, err
	}
	return v, nil
}

// BuildAggFromCreateRequestDto - build an agg.Video from dto.CreateRequest
func (b *VideoBuilder) BuildAggFromCreateRequestDto(dto dto.CreateRequest) *agg.Video {
	return &agg.Video{
		Video: entity.Video{
			Name:        dto.GetName(),
			Path:        dto.GetPath(),
			Description: dto.GetDescription(),
		},
		Timestamp: vo.Timestamp{
			CreatedAt: time.Now(),
		},
	}
}

// BuildUpdateRequestDtoFromRequest - build a dto.VideoUpdateRequestDto from raw *http.Request
func (b *VideoBuilder) BuildUpdateRequestDtoFromRequest(r *http.Request) (*dto.VideoUpdateRequestDto, error) {
	videoDto := &dto.VideoUpdateRequestDto{}
	if err := json.NewDecoder(r.Body).Decode(&videoDto); err != nil {
		return nil, err
	}

	hexId, err := b.extractor.GetParameter(id, r)
	if err != nil {
		return nil, err
	}
	oid, err := primitive.ObjectIDFromHex(hexId)
	if err != nil {
		return nil, err
	}
	videoDto.ID = vo.ID{Value: oid}

	return videoDto, nil
}

// BuildAggFromUpdateRequestDto - build an agg.Video from dto.UpdateRequest
func (b *VideoBuilder) BuildAggFromUpdateRequestDto(dto dto.UpdateRequest) (*agg.Video, error) {
	video, err := b.repository.Find(b.ctx, dto.GetId())
	if err != nil {
		return nil, err
	}

	changes := 0
	if video.Video.Name != dto.GetName() {
		video.Video.Name = dto.GetName()
		changes++
	}
	if video.Video.Description != dto.GetDescription() {
		video.Video.Description = dto.GetDescription()
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

	hexId, err := b.extractor.GetParameter(id, r)
	if err != nil {
		return nil, err
	}
	oid, err := primitive.ObjectIDFromHex(hexId)
	if err != nil {
		return nil, err
	}
	videoDto.ID = vo.ID{Value: oid}

	return videoDto, nil
}

// BuildListRequestDtoFromRequest - build a dto.ListRequest from raw *http.Request
func (b *VideoBuilder) BuildListRequestDtoFromRequest(r *http.Request) (*dto.VideoListRequestDto, error) {
	videoDto := &dto.VideoListRequestDto{}

	if b.extractor.HasParameter(name, r) {
		if nm, err := b.extractor.GetParameter(name, r); err == nil {
			videoDto.Name = nm
		}
	}
	if b.extractor.HasParameter(path, r) {
		if pth, err := b.extractor.GetParameter(path, r); err == nil {
			videoDto.Path = pth
		}
	}
	if b.extractor.HasParameter(page, r) {
		pg, _ := b.extractor.GetParameter(page, r)
		pgi, atoiErr := strconv.Atoi(pg)
		if atoiErr != nil {
			return nil, atoiErr
		}
		videoDto.Page = pgi
	} else {
		videoDto.Page = pageDefaultValue
	}
	if b.extractor.HasParameter(limit, r) {
		l, _ := b.extractor.GetParameter(limit, r)
		li, atoiErr := strconv.Atoi(l)
		if atoiErr != nil {
			return nil, atoiErr
		}
		videoDto.Limit = li
	} else {
		videoDto.Limit = limitDefaultValue
	}

	return videoDto, nil
}

func (b *VideoBuilder) BuildDeleteRequestDtoFromRequest(r *http.Request) (*dto.VideoDeleteRequestDto, error) {
	videoGetDto, err := b.BuildGetRequestDtoFromRequest(r)
	if err != nil {
		return nil, err
	}

	return &dto.VideoDeleteRequestDto{ID: videoGetDto.ID}, nil
}
