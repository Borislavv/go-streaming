package security

import (
	"errors"
	"github.com/Borislavv/video-streaming/internal/domain/errtype"
	"github.com/Borislavv/video-streaming/internal/domain/logger/interface"
	"github.com/Borislavv/video-streaming/internal/domain/service/di/interface"
	securityinterface "github.com/Borislavv/video-streaming/internal/domain/service/security/interface"
	"golang.org/x/crypto/bcrypt"
)

type PasswordHasher struct {
	logger loggerinterface.Logger
	cost   int
}

func NewPasswordHasher(serviceContainer diinterface.ServiceContainer, cost int) (*PasswordHasher, error) {
	loggerService, err := serviceContainer.GetLoggerService()
	if err != nil {
		return nil, err
	}

	return &PasswordHasher{
		logger: loggerService,
		cost:   cost,
	}, nil
}

func (s *PasswordHasher) Hash(password string) (hash string, err error) {
	hashBytes, err := bcrypt.GenerateFromPassword([]byte(password), s.cost)
	if err != nil {
		return "", s.logger.LogPropagate(err)
	}
	return string(hashBytes), nil
}

func (s *PasswordHasher) Verify(user securityinterface.Passwordness, password string) (err error) {
	if err = bcrypt.CompareHashAndPassword([]byte(user.GetPassword()), []byte(password)); err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return errtype.NewAuthFailedError("passwords did not match")
		} else {
			return s.logger.LogPropagate(err)
		}
	}
	return nil
}
