package dto

import "github.com/Borislavv/video-streaming/internal/domain/vo"

// VideoCreateRequestDto - used when u want to create a new one video
type VideoCreateRequestDto struct {
	Name        string `json:"name"`
	Path        string `json:"path"`
	Description string `json:"description,omitempty"`
}

func (req *VideoCreateRequestDto) GetName() string {
	return req.Name
}

func (req *VideoCreateRequestDto) GetPath() string {
	return req.Path
}

func (req *VideoCreateRequestDto) GetDescription() string {
	return req.Description
}

// VideoUpdateRequestDto - used when u want to update a video record
type VideoUpdateRequestDto struct {
	ID          vo.ID  `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
}

func (req *VideoUpdateRequestDto) GetId() vo.ID {
	return req.ID
}

func (req *VideoUpdateRequestDto) GetName() string {
	return req.Name
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
	Name string `json:"name"` // part of name
	Path string `json:"path"` // part of path
	PaginationRequestDto
}

func (req *VideoListRequestDto) GetName() string {
	return req.Name
}

func (req *VideoListRequestDto) GetPath() string {
	return req.Path
}

// VideoDeleteRequestDto - used when u want to remove the video
type VideoDeleteRequestDto struct {
	ID vo.ID `json:"id"`
}

func (req *VideoDeleteRequestDto) GetId() vo.ID {
	return req.ID
}
