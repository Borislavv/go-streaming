package validator

import (
	"github.com/Borislavv/video-streaming/internal/domain/dto"
	"net/http"
)

type Auth interface {
	ValidateAuthRequest(reqDTO dto.AuthRequest) error
	ValidateTokennessRequest(r *http.Request) error
}
