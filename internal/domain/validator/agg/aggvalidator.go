package aggvalidator

import "github.com/Borislavv/video-streaming/internal/domain/agg"

type Video interface {
	Validate(video *agg.Video) error
}
