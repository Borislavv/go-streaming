package authenticator

import (
	"github.com/Borislavv/video-streaming/internal/domain/dto"
	"github.com/Borislavv/video-streaming/internal/domain/dto/interface"
	"github.com/Borislavv/video-streaming/internal/domain/enum"
	"github.com/Borislavv/video-streaming/internal/domain/errors"
	"github.com/Borislavv/video-streaming/internal/domain/logger/interface"
	"github.com/Borislavv/video-streaming/internal/domain/service/di/interface"
	security_interface "github.com/Borislavv/video-streaming/internal/domain/service/security/interface"
	tokenizer_interface "github.com/Borislavv/video-streaming/internal/domain/service/tokenizer/interface"
	user_interface "github.com/Borislavv/video-streaming/internal/domain/service/user/interface"
	validator_interface "github.com/Borislavv/video-streaming/internal/domain/validator/interface"
	"github.com/Borislavv/video-streaming/internal/domain/vo"
	"net/http"
)

var (
	tokenVerificationFailed = "token verification failed"
)

type AuthService struct {
	logger         logger_interface.Logger
	userService    user_interface.CRUD
	validator      validator_interface.Auth
	tokenizer      tokenizer_interface.Tokenizer
	passwordHasher security_interface.PasswordHasher
}

func NewAuthService(serviceContainer diinterface.ContainerManager) (*AuthService, error) {
	loggerService, err := serviceContainer.GetLoggerService()
	if err != nil {
		return nil, err
	}

	userCRUDService, err := serviceContainer.GetUserCRUDService()
	if err != nil {
		return nil, loggerService.LogPropagate(err)
	}

	authValidator, err := serviceContainer.GetAuthValidator()
	if err != nil {
		return nil, loggerService.LogPropagate(err)
	}

	tokenizerService, err := serviceContainer.GetTokenizerService()
	if err != nil {
		return nil, loggerService.LogPropagate(err)
	}

	passwordHasherService, err := serviceContainer.GetPasswordHasherService()
	if err != nil {
		return nil, loggerService.LogPropagate(err)
	}

	return &AuthService{
		logger:         loggerService,
		userService:    userCRUDService,
		validator:      authValidator,
		tokenizer:      tokenizerService,
		passwordHasher: passwordHasherService,
	}, nil
}

// Auth will check raw credentials and generate a new access token for given user.
func (s *AuthService) Auth(req dto_interface.AuthRequest) (token string, err error) {
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

	// extract token from request
	token, err := s.extractToken(r)
	if err != nil {
		return vo.ID{}, s.logger.LogPropagate(err)
	}

	// validate token and extract userID from it
	userID, err = s.tokenizer.Verify(token)
	if err != nil {
		if berr := s.tokenizer.Block(token, tokenVerificationFailed); berr != nil {
			return vo.ID{}, s.logger.LogPropagate(berr)
		}
		return vo.ID{}, s.logger.LogPropagate(err)
	}

	return userID, nil
}

func (s *AuthService) extractToken(r *http.Request) (token string, err error) {
	token = r.Header.Get(enum.AccessTokenHeaderKey)
	if token != "" {
		return token, nil
	}

	cookie, err := r.Cookie(enum.AccessTokenHeaderKey)
	if err == nil {
		return cookie.Value, nil
	}

	return "", errors.NewAccessTokenIsEmptyOrOmittedError()
}
