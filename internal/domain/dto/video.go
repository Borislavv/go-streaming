package dto

import (
	"github.com/Borislavv/video-streaming/internal/domain/vo"
	"time"
)

// VideoCreateRequestDTO - used when u want to create a new one video
type VideoCreateRequestDTO struct {
	Name        string `json:"name"`
	UserID      vo.ID
	ResourceID  vo.ID  `json:"resourceID"`
	Description string `json:"description,omitempty"`
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

// VideoUpdateRequestDTO - used when u want to update a video record
type VideoUpdateRequestDTO struct {
	ID          vo.ID  `json:"id"`
	Name        string `json:"name"`
	ResourceID  vo.ID  `json:"resourceID"`
	Description string `json:"description,omitempty"`
}

func (req *VideoUpdateRequestDTO) GetID() vo.ID {
	return req.ID
}
func (req *VideoUpdateRequestDTO) GetName() string {
	return req.Name
}
func (req *VideoUpdateRequestDTO) GetResourceID() vo.ID {
	return req.ResourceID
}
func (req *VideoUpdateRequestDTO) GetDescription() string {
	return req.Description
}

// VideoGetRequestDTO - used when u want to find a single video
type VideoGetRequestDTO struct {
	ID vo.ID `json:"id"`
}

func (req *VideoGetRequestDTO) GetID() vo.ID {
	return req.ID
}

// VideoListRequestDTO - used when u want to find a collection of videos
type VideoListRequestDTO struct {
	Name      string    `json:"name"` // part of name
	CreatedAt time.Time `json:"createdAt" format:"2006-01-02T15:04:05Z07:00"`
	From      time.Time `json:"from" format:"2006-01-02T15:04:05Z07:00"`
	To        time.Time `json:"to" format:"2006-01-02T15:04:05Z07:00"`
	PaginationRequestDTO
}

func (req *VideoListRequestDTO) GetName() string {
	return req.Name
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

// VideoDeleteRequestDto - used when u want to remove the video
type VideoDeleteRequestDto struct {
	ID vo.ID `json:"id"`
}

func (req *VideoDeleteRequestDto) GetID() vo.ID {
	return req.ID
}
