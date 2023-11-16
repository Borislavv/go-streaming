package security

import "github.com/Borislavv/video-streaming/internal/domain/agg"

type PasswordHasher interface {
	Hash(password string) (hash string, err error)
	Verify(userAgg *agg.User, password string) (err error)
}
