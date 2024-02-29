package builderinterface

import (
	dtointerface "github.com/Borislavv/video-streaming/internal/domain/dto/interface"
	"net/http"
)

type Auth interface {
	BuildAuthRequestDTOFromRequest(r *http.Request) (dtointerface.AuthRequest, error)
}
