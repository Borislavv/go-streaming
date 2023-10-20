package auth

import (
	"github.com/Borislavv/video-streaming/internal/domain/dto"
	"github.com/Borislavv/video-streaming/internal/domain/errors"
	"github.com/Borislavv/video-streaming/internal/domain/logger"
	"github.com/Borislavv/video-streaming/internal/domain/service/user"
	"github.com/Borislavv/video-streaming/internal/domain/validator"
	"github.com/Borislavv/video-streaming/internal/domain/vo"
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

	userAgg, err := s.userService.Get(dto.NewUserGetRequestDTO(vo.ID{}, reqDTO.GetEmail()))
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

func (s *Service) GetToken(reqDTO dto.AuthRequest) (token string, err error) {
	// raw request validation (checking that email and pass is not empty)
	if err = s.validator.ValidateAuthRequest(reqDTO); err != nil {
		return "", s.logger.LogPropagate(err)
	}

	// getting the target user agg. by email
	userAgg, err := s.userService.Get(dto.NewUserGetRequestDTO(vo.ID{}, reqDTO.GetEmail()))
	if err != nil {
		return "", s.logger.LogPropagate(err)
	}

	// checking that credentials is valid
	if userAgg.Password != reqDTO.GetPassword() {
		return "", errors.NewAuthFailedError("passwords did not match")
	}

	// generating a new access token string
	token, err = s.tokenizer.New(userAgg)
	if err != nil {
		// TODO check if the user has a valid tokens then disable it (close all sessions)
		return "", s.logger.LogPropagate(err)
	}

	return token, nil
}
