package query

import "github.com/Borislavv/video-streaming/internal/domain/vo"

type HasBlockedToken interface {
	GetToken() string
	GetUserID() vo.ID
}
