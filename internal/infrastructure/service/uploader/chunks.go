package uploader

import "net/http"

type ChunksUploader struct {
}

func NewChunkUploader() *ChunksUploader {
	return &ChunksUploader{}
}

func (u *ChunksUploader) Upload(r *http.Request) (resourceId string, err error) {
	return "", nil
}
