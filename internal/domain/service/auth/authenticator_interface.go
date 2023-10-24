package auth

import (
	"github.com/Borislavv/video-streaming/internal/domain/dto"
)

type Authenticator interface {
	// Auth will check raw credentials and generate a new access token for given user.
	Auth(reqDTO dto.AuthRequest) (token string, err error)
}
