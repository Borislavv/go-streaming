package dto

import (
	"github.com/Borislavv/video-streaming/internal/domain/vo"
	"time"
)

// VideoCreateRequestDTO - used when u want to create a new one video.
type VideoCreateRequestDTO struct {
	/*Required*/ Name string `json:"name"`
	/*Required*/ UserID vo.ID
	/*Required*/ ResourceID vo.ID `json:"resourceID"`
	/*Optional*/ Description string `json:"description,omitempty"`
}

func (req *VideoCreateRequestDTO) GetName() string {
	return req.Name
}
func (req *VideoCreateRequestDTO) GetUserID() vo.ID {
	return req.UserID
}
func (req *VideoCreateRequestDTO) GetResourceID() vo.ID {
	return req.ResourceID
}
func (req *VideoCreateRequestDTO) GetDescription() string {
	return req.Description
}

// VideoUpdateRequestDTO - used when u want to update a video record.
type VideoUpdateRequestDTO struct {
	/*Required*/ ID vo.ID `json:"id"`
	/*Optional*/ Name string `json:"name"`
	/*Optional*/ UserID vo.ID
	/*Optional*/ ResourceID vo.ID `json:"resourceID"`
	/*Optional*/ Description string `json:"description,omitempty"`
}

func (req *VideoUpdateRequestDTO) GetID() vo.ID {
	return req.ID
}
func (req *VideoUpdateRequestDTO) GetName() string {
	return req.Name
}
func (req *VideoUpdateRequestDTO) GetUserID() vo.ID {
	return req.UserID
}
func (req *VideoUpdateRequestDTO) GetResourceID() vo.ID {
	return req.ResourceID
}
func (req *VideoUpdateRequestDTO) GetDescription() string {
	return req.Description
}

// VideoGetRequestDTO - used when you want to find a single video by Name or ID, but you always must specify a UserID.
type VideoGetRequestDTO struct {
	/*Optional*/ ID vo.ID `json:"id"`
	/*Optional*/ Name string
	/*Required*/ UserID vo.ID
}

func NewVideoGetRequestDTO(id vo.ID, name string, userID vo.ID) *VideoGetRequestDTO {
	return &VideoGetRequestDTO{
		ID:     id,
		Name:   name,
		UserID: userID,
	}
}
func (req *VideoGetRequestDTO) GetID() vo.ID {
	return req.ID
}
func (req *VideoGetRequestDTO) GetName() string {
	return req.Name
}
func (req *VideoGetRequestDTO) GetUserID() vo.ID {
	return req.UserID
}

// VideoListRequestDTO - used when u want to find a collection of videos.
type VideoListRequestDTO struct {
	/*Required*/ UserID vo.ID
	/*Optional*/ Name string `json:"name"` // part of name
	/*Optional*/ CreatedAt time.Time `json:"createdAt" format:"2006-01-02T15:04:05Z07:00"`
	/*Optional*/ From time.Time `json:"from" format:"2006-01-02T15:04:05Z07:00"`
	/*Optional*/ To time.Time `json:"to" format:"2006-01-02T15:04:05Z07:00"`
	/*Optional*/ PaginationRequestDTO
}

func NewVideoListRequestDTO(
	name string, userID vo.ID,
	createdAt time.Time, from time.Time, to time.Time,
	page int, limit int,
) *VideoListRequestDTO {
	return &VideoListRequestDTO{
		Name:      name,
		UserID:    userID,
		CreatedAt: createdAt,
		From:      from,
		To:        to,
		PaginationRequestDTO: PaginationRequestDTO{
			Page:  page,
			Limit: limit,
		},
	}
}
func (req *VideoListRequestDTO) GetName() string {
	return req.Name
}
func (req *VideoListRequestDTO) GetUserID() vo.ID {
	return req.UserID
}
func (req *VideoListRequestDTO) GetCreatedAt() time.Time {
	return req.CreatedAt
}
func (req *VideoListRequestDTO) GetFrom() time.Time {
	return req.From
}
func (req *VideoListRequestDTO) GetTo() time.Time {
	return req.To
}

// VideoDeleteRequestDto - used when you want to remove the video.
type VideoDeleteRequestDto struct {
	/*Required*/ ID vo.ID `json:"id"`
	/*Required*/ UserID vo.ID
}

func NewVideoDeleteRequestDto(id vo.ID, userID vo.ID) *VideoDeleteRequestDto {
	return &VideoDeleteRequestDto{
		ID:     id,
		UserID: userID,
	}
}
func (req *VideoDeleteRequestDto) GetID() vo.ID {
	return req.ID
}
func (req *VideoDeleteRequestDto) GetUserID() vo.ID {
	return req.UserID
}
