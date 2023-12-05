package security

import (
	"github.com/Borislavv/video-streaming/internal/domain/errors"
	"github.com/Borislavv/video-streaming/internal/domain/logger"
	"github.com/Borislavv/video-streaming/internal/domain/service/security"
	"golang.org/x/crypto/bcrypt"
)

type PasswordHasher struct {
	logger logger.Logger
	cost   int
}

func NewPasswordHasher(logger logger.Logger, cost int) *PasswordHasher {
	return &PasswordHasher{
		logger: logger,
		cost:   cost,
	}
}

func (s *PasswordHasher) Hash(password string) (hash string, err error) {
	hashBytes, err := bcrypt.GenerateFromPassword([]byte(password), s.cost)
	if err != nil {
		return "", s.logger.LogPropagate(err)
	}
	return string(hashBytes), nil
}

func (s *PasswordHasher) Verify(user security.Passwordness, password string) (err error) {
	if err = bcrypt.CompareHashAndPassword([]byte(user.GetPassword()), []byte(password)); err != nil {
		if err == bcrypt.ErrMismatchedHashAndPassword {
			return errors.NewAuthFailedError("passwords did not match")
		} else {
			return s.logger.LogPropagate(err)
		}
	}
	return nil
}
