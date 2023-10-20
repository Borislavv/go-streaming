package auth

import (
	"github.com/Borislavv/video-streaming/internal/domain/dto"
	"github.com/Borislavv/video-streaming/internal/domain/errors"
	"github.com/Borislavv/video-streaming/internal/domain/logger"
	"github.com/Borislavv/video-streaming/internal/domain/service/user"
	"github.com/Borislavv/video-streaming/internal/domain/validator"
	"net/http"
)

type Service struct {
	logger      logger.Logger
	userService user.CRUD
	validator   validator.Auth
	tokenizer   Tokenizer
}

func NewService(
	logger logger.Logger,
	userService user.CRUD,
	validator validator.Auth,
	tokenizer Tokenizer,
) *Service {
	return &Service{
		logger:      logger,
		userService: userService,
		validator:   validator,
		tokenizer:   tokenizer,
	}
}

func (s *Service) Auth(w http.ResponseWriter, r *http.Request, reqDTO dto.AuthRequest) (token string, err error) {
	if err = s.validator.ValidateAuthRequest(reqDTO); err != nil {
		return "", s.logger.LogPropagate(err)
	}

	userAgg, err := s.userService.Get(&dto.UserGetRequestDTO{Email: reqDTO.GetEmail()})
	if err != nil {
		return "", s.logger.LogPropagate(err)
	}

	if userAgg.Password != reqDTO.GetPassword() {
		return "", errors.NewAuthFailedError("passwords did not match")
	}

	if err = s.tokenizer.Set(w, r, userAgg); err != nil {
		return "", s.logger.LogPropagate(err)
	}

	token, err = s.tokenizer.Get(r)
	if err != nil {
		return "", s.logger.LogPropagate(err)
	}

	return token, nil
}
