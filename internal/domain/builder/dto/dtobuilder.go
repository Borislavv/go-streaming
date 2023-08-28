package dtobuilder

import (
	"github.com/Borislavv/video-streaming/internal/domain/dto"
	"net/http"
)

type Video interface {
	BuildFromRequest(r *http.Request) (*dto.Video, error)
}
