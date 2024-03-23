package detectorinterface

import "github.com/Borislavv/video-streaming/internal/domain/entity"

type Codecs interface {
	Detect(resource entity.Resource) (audioCodec string, videoCodec string, err error)
}
