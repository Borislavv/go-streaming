package validator_interface

import (
	dto_interface "github.com/Borislavv/video-streaming/internal/domain/dto/interface"
	"net/http"
)

type Auth interface {
	ValidateAuthRequest(reqDTO dto_interface.AuthRequest) error
	ValidateTokennessRequest(r *http.Request) error
}
