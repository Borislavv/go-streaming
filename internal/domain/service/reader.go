package service

import (
	"github.com/Borislavv/video-streaming/internal/infrastructure/service/reader"
	"github.com/Borislavv/video-streaming/internal/infrastructure/service/streamer"
)

type Reader interface {
	Read(resource reader.Resource) chan streamer.Chunk
}
