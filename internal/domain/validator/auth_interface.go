package validator

import "github.com/Borislavv/video-streaming/internal/domain/dto"

type Auth interface {
	ValidateAuthRequest(reqDTO dto.AuthRequest) error
}
