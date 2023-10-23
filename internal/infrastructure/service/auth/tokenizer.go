package auth

import (
	"github.com/Borislavv/video-streaming/internal/domain/agg"
	"github.com/Borislavv/video-streaming/internal/domain/errors"
	"github.com/Borislavv/video-streaming/internal/domain/logger"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
	"time"
)

const TokenCookieKey = "access-token"

type TokenizerService struct {
	logger                  logger.Logger
	jwtTokenAcceptedIssuers []string
	jwtTokenIssuer          string
	jwtSecretSalt           string
	jwtTokenEncryptAlgo     string
	jwtTokenExpiresAfter    int64
}

func NewTokenizerService(
	logger logger.Logger,
	jwtSecretSalt string,
) *TokenizerService {
	return &TokenizerService{
		logger:        logger,
		jwtSecretSalt: jwtSecretSalt,
	}
}

func (s *TokenizerService) New(user *agg.User) (token string, err error) {
	tkn := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": user.ID,
		"iss": s.jwtTokenIssuer,
		"exp": &jwt.NumericDate{Time: time.Now().Add(time.Second * time.Duration(s.jwtTokenExpiresAfter))},
	})

	token, err = tkn.SignedString([]byte(s.jwtSecretSalt))
	if err != nil {
		return "", s.logger.LogPropagate(err)
	}
	return token, nil
}

func (s *TokenizerService) Set(w http.ResponseWriter, r *http.Request, user *agg.User) error {
	if !s.Has(r) {
		if err := s.Refresh(w, user); err != nil {
			return s.logger.LogPropagate(err)
		}
	}
	return nil
}

func (s *TokenizerService) Has(r *http.Request) bool {
	cookie, err := r.Cookie(TokenCookieKey)
	if err != nil {
		if err != http.ErrNoCookie {
			s.logger.Log(err)
		}
		return false
	}
	return cookie.Value != ""
}

func (s *TokenizerService) Get(r *http.Request) (token string, err error) {
	cookie, err := r.Cookie(TokenCookieKey)
	if err != nil {
		return "", s.logger.LogPropagate(err)
	}
	return cookie.Value, nil
}

func (s *TokenizerService) Refresh(w http.ResponseWriter, user *agg.User) error {
	token, err := s.New(user)
	if err != nil {
		return s.logger.LogPropagate(err)
	}

	cookie := &http.Cookie{
		Name:    TokenCookieKey,
		Value:   token,
		Expires: time.Now().Add(time.Second * time.Duration(s.jwtTokenExpiresAfter)),
	}

	http.SetCookie(w, cookie)

	return nil
}

func (s *TokenizerService) IsValid(w http.ResponseWriter, r *http.Request, user *agg.User) (ok bool, err error) {
	if !s.Has(r) {
		return false, nil
	}

	if ok, err = s.isValid(r, user); err != nil || !ok {
		s.Remove(w)
	}

	return ok, err
}

func (s *TokenizerService) isValid(r *http.Request, user *agg.User) (ok bool, err error) {
	// extract token string from request
	givenToken, err := s.Get(r)
	if err != nil {
		return false, s.logger.LogPropagate(err)
	}

	parsedToken, err := jwt.Parse(givenToken, func(token *jwt.Token) (interface{}, error) {
		if token.Header["alg"] != s.jwtTokenEncryptAlgo {
			// user must be banned here because the algo wasn't matched
			return nil, errors.NewTokenAlgoWasNotMatchedError()
		}
		// cast to the configured token signature type (stored in `s.jwtTokenEncryptAlgo`)
		if _, success := token.Method.(*jwt.SigningMethodHMAC); !success {
			return nil, errors.NewTokenUnexpectedSigningMethodError(token.Header["alg"])
		}
		// jwtSecretSalt is a string containing your secret, but you need pass the []byte
		return []byte(s.jwtSecretSalt), nil
	})
	if err != nil {
		// parsing token error occurred
		return false, s.logger.LogPropagate(err)
	}

	// extracting claims of the token payload
	if claims, success := parsedToken.Claims.(jwt.MapClaims); success && parsedToken.Valid {
		// extracting subject from the claims
		sub, err := claims.GetSubject()
		if err != nil {
			return false, s.logger.LogPropagate(err)
		}
		// checking the subject is equals to current user
		if sub != user.ID.Value.Hex() {
			return false, errors.NewTokenSubjectPayloadWasNotMatchedError(user.ID.Value.Hex(), sub)
		}
	} else {
		// error occurred while extracting claims from token or token is not valid
		return false, errors.NewTokenInvalidError(givenToken)
	}

	// all checks were passed, token is valid
	return true, nil
}

func (s *TokenizerService) Remove(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name:    TokenCookieKey,
		Value:   "",
		Expires: time.Unix(0, 0),
	})
}
