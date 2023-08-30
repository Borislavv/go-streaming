package dtovalidator

import "github.com/Borislavv/video-streaming/internal/domain/dto"

type Video interface {
	Validate(video *dto.Video) error
}
