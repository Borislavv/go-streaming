package uploader

import (
	"github.com/Borislavv/video-streaming/internal/domain/dto"
)

type ChunkUploader struct {
}

func NewChunkUploader() *ChunkUploader {
	return &ChunkUploader{}
}

func (u *ChunkUploader) Upload(req dto.UploadRequest) (err error) {
	// todo must be implemented
	return nil
}
