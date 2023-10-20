package auth

import (
	"github.com/Borislavv/video-streaming/internal/domain/dto"
	"net/http"
)

type Authenticator interface {
	Auth(w http.ResponseWriter, r *http.Request, reqDTO dto.AuthRequest) (token string, err error)
	GetToken(reqDTO dto.AuthRequest) (token string, err error)
}
