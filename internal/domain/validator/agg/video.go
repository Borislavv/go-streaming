package aggvalidator

import (
	"errors"
	"github.com/Borislavv/video-streaming/internal/domain/agg"
)

type VideoAggValidator struct {
}

func NewVideoAggValidator() *VideoAggValidator {
	return &VideoAggValidator{}
}

func (v *VideoAggValidator) Validate(video *agg.Video) error {
	if video.Video.Name == "" {
		return errors.New("video.name must not be empty")
	} else if video.Video.Path == "" {
		return errors.New("video.path must not be empty")
	}
	return nil
}
