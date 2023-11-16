package security

import (
	"github.com/Borislavv/video-streaming/internal/domain/agg"
	"github.com/Borislavv/video-streaming/internal/domain/logger"
	"github.com/Borislavv/video-streaming/internal/domain/repository"
	"golang.org/x/crypto/bcrypt"
)

type PasswordHasher struct {
	logger     logger.Logger
	repository repository.User
	cost       int
}

func NewPasswordHasher(logger logger.Logger, repository repository.User, cost int) *PasswordHasher {
	return &PasswordHasher{
		logger:     logger,
		repository: repository,
		cost:       cost,
	}
}

func (s *PasswordHasher) Hash(password string) (hash string, err error) {
	hashBytes, err := bcrypt.GenerateFromPassword([]byte(password), s.cost)
	if err != nil {
		return "", s.logger.LogPropagate(err)
	}
	return string(hashBytes), nil
}

func (s *PasswordHasher) Verify(userAgg *agg.User, password string) (err error) {
	return bcrypt.CompareHashAndPassword([]byte(userAgg.Password), []byte(password))
}
