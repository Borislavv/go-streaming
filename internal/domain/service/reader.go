package service

import (
	"github.com/Borislavv/video-streaming/internal/domain/dto"
)

type Reader interface {
	Read(resource dto.Resource) chan dto.Chunk
}
