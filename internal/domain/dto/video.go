package dto

import (
	"github.com/Borislavv/video-streaming/internal/domain/vo"
	"time"
)

// VideoCreateRequestDto - used when u want to create a new one video
type VideoCreateRequestDto struct {
	Name        string `json:"name"`
	ResourceID  vo.ID  `json:"resourceID"`
	Description string `json:"description,omitempty"`
}

func (req *VideoCreateRequestDto) GetName() string {
	return req.Name
}
func (req *VideoCreateRequestDto) GetResourceID() vo.ID {
	return req.ResourceID
}
func (req *VideoCreateRequestDto) GetDescription() string {
	return req.Description
}

// VideoUpdateRequestDto - used when u want to update a video record
type VideoUpdateRequestDto struct {
	ID          vo.ID  `json:"id"`
	Name        string `json:"name"`
	ResourceID  vo.ID  `json:"resourceID"`
	Description string `json:"description,omitempty"`
}

func (req *VideoUpdateRequestDto) GetId() vo.ID {
	return req.ID
}
func (req *VideoUpdateRequestDto) GetName() string {
	return req.Name
}
func (req *VideoUpdateRequestDto) GetResourceID() vo.ID {
	return req.ResourceID
}
func (req *VideoUpdateRequestDto) GetDescription() string {
	return req.Description
}

// VideoGetRequestDto - used when u want to find a single video
type VideoGetRequestDto struct {
	ID vo.ID `json:"id"`
}

func (req *VideoGetRequestDto) GetId() vo.ID {
	return req.ID
}

// VideoListRequestDto - used when u want to find a collection of videos
type VideoListRequestDto struct {
	Name      string    `json:"name"` // part of name
	CreatedAt time.Time `json:"createdAt" format:"2006-01-02T15:04:05Z07:00"`
	From      time.Time `json:"from" format:"2006-01-02T15:04:05Z07:00"`
	To        time.Time `json:"to" format:"2006-01-02T15:04:05Z07:00"`
	PaginationRequestDto
}

func (req *VideoListRequestDto) GetName() string {
	return req.Name
}
func (req *VideoListRequestDto) GetCreatedAt() time.Time {
	return req.CreatedAt
}
func (req *VideoListRequestDto) GetFrom() time.Time {
	return req.From
}
func (req *VideoListRequestDto) GetTo() time.Time {
	return req.To
}

// VideoDeleteRequestDto - used when u want to remove the video
type VideoDeleteRequestDto struct {
	ID vo.ID `json:"id"`
}

func (req *VideoDeleteRequestDto) GetId() vo.ID {
	return req.ID
}
