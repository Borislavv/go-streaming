package _interface

import (
	"github.com/Borislavv/video-streaming/internal/domain/dto"
	"net/http"
)

type Auth interface {
	BuildAuthRequestDTOFromRequest(r *http.Request) (dto.AuthRequest, error)
}
