package dto

import "github.com/Borislavv/video-streaming/internal/domain/vo"

type VideoCreateRequestDto struct {
	Name        string `json:"name"`
	Path        string `json:"path"`
	Description string `json:"description,omitempty"`
}

func (dto *VideoCreateRequestDto) GetName() string {
	return dto.Name
}

func (dto *VideoCreateRequestDto) GetPath() string {
	return dto.Path
}

func (dto *VideoCreateRequestDto) GetDescription() string {
	return dto.Description
}

type VideoUpdateRequestDto struct {
	ID          vo.ID  `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
}

func (dto *VideoUpdateRequestDto) GetId() vo.ID {
	return dto.ID
}

func (dto *VideoUpdateRequestDto) GetName() string {
	return dto.Name
}

func (dto *VideoUpdateRequestDto) GetDescription() string {
	return dto.Description
}

type VideoListRequestDto struct {
	Name       string               `json:"name"` // part of name
	Path       string               `json:"path"` // part of path
	Pagination PaginationRequestDto `json:"pagination"`
}

func (dto *VideoListRequestDto) GetName() string {
	return dto.Name
}

func (dto *VideoListRequestDto) GetPath() string {
	return dto.Path
}
