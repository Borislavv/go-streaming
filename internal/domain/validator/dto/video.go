package dtovalidator

import (
	"errors"
	"github.com/Borislavv/video-streaming/internal/domain/dto"
)

type VideoDtoValidator struct {
}

func NewVideoDtoValidator() *VideoDtoValidator {
	return &VideoDtoValidator{}
}

func (v *VideoDtoValidator) Validate(video *dto.Video) error {
	if video.Name == "" {
		return errors.New("video.name must not be empty")
	} else if video.Path == "" {
		return errors.New("video.path must not be empty")
	}
	return nil
}
