package validatorinterface

import (
	dtointerface "github.com/Borislavv/video-streaming/internal/domain/dto/interface"
	"net/http"
)

type Auth interface {
	ValidateAuthRequest(reqDTO dtointerface.AuthRequest) error
	ValidateTokennessRequest(r *http.Request) error
}
