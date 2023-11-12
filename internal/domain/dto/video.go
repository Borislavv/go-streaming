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
	UserID      vo.ID
	ResourceID  vo.ID  `json:"resourceID"`
	Description string `json:"description,omitempty"`
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

// VideoGetRequestDTO - used when u want to find a single video
type VideoGetRequestDTO struct {
	ID     vo.ID `json:"id"`
	UserID vo.ID
}

func NewVideoGetRequestDTO(id vo.ID, userID vo.ID) *VideoGetRequestDTO {
	return &VideoGetRequestDTO{
		ID:     id,
		UserID: userID,
	}
}
func (req *VideoGetRequestDTO) GetID() vo.ID {
	return req.ID
}
func (req *VideoGetRequestDTO) GetUserID() vo.ID {
	return req.UserID
}

// VideoListRequestDTO - used when u want to find a collection of videos
type VideoListRequestDTO struct {
	Name      string `json:"name"` // part of name
	UserID    vo.ID
	CreatedAt time.Time `json:"createdAt" format:"2006-01-02T15:04:05Z07:00"`
	From      time.Time `json:"from" format:"2006-01-02T15:04:05Z07:00"`
	To        time.Time `json:"to" format:"2006-01-02T15:04:05Z07:00"`
	PaginationRequestDTO
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

// VideoDeleteRequestDto - used when u want to remove the video
type VideoDeleteRequestDto struct {
	ID     vo.ID `json:"id"`
	UserID vo.ID
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
