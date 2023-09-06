package builder

import (
	"context"
	"encoding/json"
	"github.com/Borislavv/video-streaming/internal/domain/agg"
	"github.com/Borislavv/video-streaming/internal/domain/dto"
	"github.com/Borislavv/video-streaming/internal/domain/entity"
	"github.com/Borislavv/video-streaming/internal/domain/repository"
	"github.com/Borislavv/video-streaming/internal/domain/vo"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
	"time"
)

type VideoBuilder struct {
	ctx        context.Context
	repository repository.Video
}

func NewVideoBuilder(ctx context.Context, repository repository.Video) *VideoBuilder {
	return &VideoBuilder{ctx: ctx, repository: repository}
}

// BuildCreateRequestDtoFromRequest - build a dto.VideoCreateRequestDto from raw http.Request
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

// BuildUpdateRequestDtoFromRequest - build a dto.VideoUpdateRequestDto from raw http.Request
func (b *VideoBuilder) BuildUpdateRequestDtoFromRequest(r *http.Request) (*dto.VideoUpdateRequestDto, error) {
	videoDto := &dto.VideoUpdateRequestDto{}
	if err := json.NewDecoder(r.Body).Decode(&videoDto); err != nil {
		return nil, err
	}

	hexId, found := mux.Vars(r)["id"]
	if found {
		oid, err := primitive.ObjectIDFromHex(hexId)
		if err != nil {
			return nil, err
		}
		videoDto.ID = vo.ID{Value: oid}
	}

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

func (b *VideoBuilder) BuildGetRequestDtoFromRequest(r *http.Request) (*dto.VideoGetRequestDto, error) {
	videoDto := &dto.VideoGetRequestDto{}

	hexId, found := mux.Vars(r)["id"]
	if found {
		oid, err := primitive.ObjectIDFromHex(hexId)
		if err != nil {
			return nil, err
		}
		videoDto.ID = vo.ID{Value: oid}
	}
	return videoDto, nil
}

func (b *VideoBuilder) BuildListRequestDtoFromRequest(r *http.Request) (*dto.VideoListRequestDto, error) {
	v := &dto.VideoListRequestDto{}
	if err := json.NewDecoder(r.Body).Decode(&v); err != nil {
		return nil, err
	}
	return v, nil
}
