package auth

import (
	"github.com/Borislavv/video-streaming/internal/domain/dto"
)

type Authenticator interface {
	AuthRaw(reqDTO dto.AuthRequest) (token string, err error)
}
