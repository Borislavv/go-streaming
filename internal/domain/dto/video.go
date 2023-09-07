package dto

import "github.com/Borislavv/video-streaming/internal/domain/vo"

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

type VideoDeleteRequestDto VideoGetRequestDto
type VideoGetRequestDto struct {
	ID vo.ID `json:"id"`
}

func (req *VideoGetRequestDto) GetId() vo.ID {
	return req.ID
}

type VideoListRequestDto struct {
	Name                 string `json:"name"` // part of name
	Path                 string `json:"path"` // part of path
	PaginationRequestDto `json:"pagination"`
}

func (req *VideoListRequestDto) GetName() string {
	return req.Name
}

func (req *VideoListRequestDto) GetPath() string {
	return req.Path
}
