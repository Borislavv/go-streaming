package detector

import (
	"context"
	"github.com/Borislavv/video-streaming/internal/domain/entity"
	"github.com/Borislavv/video-streaming/internal/domain/logger/interface"
	"github.com/Borislavv/video-streaming/internal/domain/service/di/interface"
	"gopkg.in/vansante/go-ffprobe.v2"
	"os"
)

type ResourceCodecs struct {
	ctx    context.Context
	logger loggerinterface.Logger
}

func NewResourceCodecs(serviceContainer diinterface.ServiceContainer) (*ResourceCodecs, error) {
	loggerService, err := serviceContainer.GetLoggerService()
	if err != nil {
		return nil, err
	}

	ctx, err := serviceContainer.GetCtx()
	if err != nil {
		return nil, loggerService.LogPropagate(err)
	}

	return &ResourceCodecs{
		ctx:    ctx,
		logger: loggerService,
	}, nil
}

// Detect will determine video and audio stream codecs of target resource
func (d *ResourceCodecs) Detect(
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
