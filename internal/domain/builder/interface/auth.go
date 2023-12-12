package builder_interface

import (
	"github.com/Borislavv/video-streaming/internal/domain/dto/interface"
	"net/http"
)

type Auth interface {
	BuildAuthRequestDTOFromRequest(r *http.Request) (dto_interface.AuthRequest, error)
}
