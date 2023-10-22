package detector

import "github.com/Borislavv/video-streaming/internal/domain/entity"

type Detector interface {
	Detect(resource entity.Resource) (audioCodec string, videoCodec string, err error)
}
