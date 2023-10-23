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

type AuthenticatorService struct {
	logger      logger.Logger
	userService user.CRUD
	validator   validator.Auth
	tokenizer   Tokenizer
}

func NewAuthenticatorService(
	logger logger.Logger,
	userService user.CRUD,
	validator validator.Auth,
	tokenizer Tokenizer,
) *AuthenticatorService {
	return &AuthenticatorService{
		logger:      logger,
		userService: userService,
		validator:   validator,
		tokenizer:   tokenizer,
	}
}

// GetToken will check credentials and generate a new access token for given user.
func (s *AuthenticatorService) GetToken(reqDTO dto.AuthRequest) (token string, err error) {
	// raw request validation (checking that email and pass is not empty)
	if err = s.validator.ValidateAuthRequest(reqDTO); err != nil {
		return "", s.logger.LogPropagate(err)
	}

	// getting the target user agg. by email
	userAgg, err := s.userService.Get(dto.NewUserGetRequestDTO(vo.ID{}, reqDTO.GetEmail()))
	if err != nil {
		return "", s.logger.LogPropagate(err)
	}

	// checking that credentials are valid
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

// SetCookie will check credentials, generate a new token and set it up in the cookies.
func (s *AuthenticatorService) SetCookie(w http.ResponseWriter, r *http.Request, reqDTO dto.AuthRequest) error {
	// raw request validation (checking that email and pass is not empty)
	if err := s.validator.ValidateAuthRequest(reqDTO); err != nil {
		return s.logger.LogPropagate(err)
	}

	// getting the target user agg. by email
	userAgg, err := s.userService.Get(dto.NewUserGetRequestDTO(vo.ID{}, reqDTO.GetEmail()))
	if err != nil {
		return s.logger.LogPropagate(err)
	}

	// checking that credentials are valid
	if userAgg.Password != reqDTO.GetPassword() {
		return errors.NewAuthFailedError("passwords did not match")
	}

	// setting up the access token in cookies
	if err = s.tokenizer.Set(w, r, userAgg); err != nil {
		return s.logger.LogPropagate(err)
	}

	return nil
}
