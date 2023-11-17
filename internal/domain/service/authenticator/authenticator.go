package authenticator

import (
	"github.com/Borislavv/video-streaming/internal/domain/dto"
	"github.com/Borislavv/video-streaming/internal/domain/enum"
	"github.com/Borislavv/video-streaming/internal/domain/errors"
	"github.com/Borislavv/video-streaming/internal/domain/logger"
	"github.com/Borislavv/video-streaming/internal/domain/service/security"
	"github.com/Borislavv/video-streaming/internal/domain/service/tokenizer"
	"github.com/Borislavv/video-streaming/internal/domain/service/user"
	"github.com/Borislavv/video-streaming/internal/domain/validator"
	"github.com/Borislavv/video-streaming/internal/domain/vo"
	"net/http"
)

type AuthService struct {
	logger         logger.Logger
	userService    user.CRUD
	validator      validator.Auth
	tokenizer      tokenizer.Tokenizer
	passwordHasher security.PasswordHasher
}

func NewAuthService(
	logger logger.Logger,
	userService user.CRUD,
	validator validator.Auth,
	tokenizer tokenizer.Tokenizer,
	passwordHasher security.PasswordHasher,
) *AuthService {
	return &AuthService{
		logger:         logger,
		userService:    userService,
		validator:      validator,
		tokenizer:      tokenizer,
		passwordHasher: passwordHasher,
	}
}

// Auth will check raw credentials and generate a new access token for given user.
func (s *AuthService) Auth(req dto.AuthRequest) (token string, err error) {
	// raw request validation (checking that email and pass is not empty)
	if err = s.validator.ValidateAuthRequest(req); err != nil {
		return "", s.logger.LogPropagate(err)
	}

	// getting the target user agg. by email
	userAgg, err := s.userService.Get(dto.NewUserGetRequestDTO(vo.ID{}, req.GetEmail()))
	if err != nil {
		return "", s.logger.LogPropagate(err)
	}

	// checking that credentials are valid
	if err = s.passwordHasher.Verify(userAgg, req.GetPassword()); err != nil {
		return "", s.logger.LogPropagate(err)
	}

	// generating a new access token string
	token, err = s.tokenizer.New(userAgg)
	if err != nil {
		// TODO check if the user has a valid tokens then disable it (close all sessions)
		return "", s.logger.LogPropagate(err)
	}

	return token, nil
}

// IsAuthed with check that token is valid and extract userID from it.
func (s *AuthService) IsAuthed(r *http.Request) (userID vo.ID, err error) {
	// validate that token is present into request headers
	if err = s.validator.ValidateTokennessRequest(r); err != nil {
		return vo.ID{}, s.logger.LogPropagate(err)
	}

	// extract token from request headers
	token := r.Header.Get(enum.AccessTokenHeaderKey)
	if token == "" {
		return vo.ID{}, s.logger.LogPropagate(errors.NewAccessTokenIsEmptyOrOmittedError())
	}

	// validate token and extract userID from it
	userID, err = s.tokenizer.Validate(token)
	if err != nil {
		return vo.ID{}, s.logger.LogPropagate(err)
	}

	return userID, nil
}
