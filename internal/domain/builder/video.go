package builder

import (
	"encoding/json"
	"github.com/Borislavv/video-streaming/internal/domain/agg"
	"github.com/Borislavv/video-streaming/internal/domain/dto"
	"github.com/Borislavv/video-streaming/internal/domain/entity"
	"github.com/Borislavv/video-streaming/internal/domain/repository"
	"github.com/Borislavv/video-streaming/internal/domain/vo"
	"net/http"
	"time"
)

type VideoBuilder struct {
	repository repository.Video
}

func NewVideoBuilder(repository repository.Video) *VideoBuilder {
	return &VideoBuilder{repository: repository}
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
	v := &dto.VideoUpdateRequestDto{}

	if err := json.NewDecoder(r.Body).Decode(&v); err != nil {
		return nil, err
	}

	return v, nil
}

// BuildAggFromUpdateRequestDto - build an agg.Video from dto.UpdateRequest
func (b *VideoBuilder) BuildAggFromUpdateRequestDto(dto dto.UpdateRequest) (*agg.Video, error) {
	//videoAgg, err := b.repository.Find(dto.GetId())
	//if err != nil {
	//	return nil, err
	//}
	//
	//if videoAgg.Name != dto.GetName() {
	//	videoAgg.Name = dto.GetName()
	//}
	//if videoAgg.Description != dto.GetDescription() {
	//	videoAgg.Description = dto.GetDescription()
	//}
	//
	//return videoAgg, nil
	return nil, nil
}

func (b *VideoBuilder) BuildListRequestDtoFromRequest(r *http.Request) (*dto.VideoListRequestDto, error) {
	v := &dto.VideoListRequestDto{}

	if err := json.NewDecoder(r.Body).Decode(&v); err != nil {
		return nil, err
	}

	return v, nil
}
