package tokenizer_interface

import (
	"github.com/Borislavv/video-streaming/internal/domain/agg"
	"github.com/Borislavv/video-streaming/internal/domain/vo"
)

type Tokenizer interface {
	New(user *agg.User) (token string, err error)
	Verify(token string) (userID vo.ID, err error)
	Block(token string, reason string) error
}
