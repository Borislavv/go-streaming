package authenticatorinterface

import (
	dtointerface "github.com/Borislavv/video-streaming/internal/domain/dto/interface"
	"github.com/Borislavv/video-streaming/internal/domain/vo"
	"net/http"
)

type Authenticator interface {
	// Auth will check raw credentials and generate a new access token for given user.
	Auth(reqDTO dtointerface.AuthRequest) (token string, err error)
	// IsAuthed with check that token is valid and extract userID from it.
	IsAuthed(r *http.Request) (userID vo.ID, err error)
}
