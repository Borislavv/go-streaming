package detector

import (
	"context"
	"github.com/Borislavv/video-streaming/internal/domain/entity"
	"github.com/Borislavv/video-streaming/internal/domain/logger"
	"gopkg.in/vansante/go-ffprobe.v2"
	"os"
)

type ResourceCodecDetector struct {
	ctx    context.Context
	logger logger.Logger
}

func NewResourceCodecInfo(ctx context.Context, logger logger.Logger) *ResourceCodecDetector {
	return &ResourceCodecDetector{
		ctx:    ctx,
		logger: logger,
	}
}

// Detect will determine video and audio stream codecs of target resource
func (d *ResourceCodecDetector) Detect(
	resource entity.Resource,
) (
	audioCodec string,
	videoCodec string,
	e error,
) {
	file, err := os.Open(resource.GetFilepath())
	if err != nil {
		return "", "", d.logger.LogPropagate(err)
	}
	defer func() { _ = file.Close() }()

	data, err := ffprobe.ProbeReader(d.ctx, file)
	if err != nil {
		return "", "", d.logger.LogPropagate(err)
	}

	audioCodec = ""
	videoCodec = ""
	if data.FirstAudioStream() != nil {
		audioCodec = data.FirstAudioStream().CodecTagString
	}
	if data.FirstVideoStream() != nil {
		videoCodec = data.FirstVideoStream().CodecTagString
	}

	return audioCodec, videoCodec, nil
}
