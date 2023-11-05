package tokenizer

import (
	"context"
	"github.com/Borislavv/video-streaming/internal/domain/agg"
	"github.com/Borislavv/video-streaming/internal/domain/errors"
	"github.com/Borislavv/video-streaming/internal/domain/logger"
	"github.com/Borislavv/video-streaming/internal/domain/repository"
	"github.com/Borislavv/video-streaming/internal/domain/vo"
	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type JwtService struct {
	ctx                     context.Context
	logger                  logger.Logger
	blockedTokenRepository  repository.BlockedToken
	jwtTokenAcceptedIssuers []string
	jwtSecretSalt           []byte
	jwtTokenIssuer          string
	jwtTokenEncryptAlgo     string
	jwtTokenExpiresAfter    int64
}

func NewJwtService(
	ctx context.Context,
	logger logger.Logger,
	blockedTokenRepository repository.BlockedToken,
	jwtTokenAcceptedIssuers []string,
	jwtSecretSalt string,
	jwtTokenIssuer string,
	jwtTokenEncryptAlgo string,
	jwtTokenExpiresAfter int64,
) *JwtService {
	return &JwtService{
		ctx:                     ctx,
		logger:                  logger,
		blockedTokenRepository:  blockedTokenRepository,
		jwtTokenAcceptedIssuers: jwtTokenAcceptedIssuers,
		jwtSecretSalt:           []byte(jwtSecretSalt),
		jwtTokenIssuer:          jwtTokenIssuer,
		jwtTokenEncryptAlgo:     jwtTokenEncryptAlgo,
		jwtTokenExpiresAfter:    jwtTokenExpiresAfter,
	}
}

// New will generate a new JWT.
func (s *JwtService) New(user *agg.User) (token string, err error) {
	tkn := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": user.ID.Value.Hex(),
		"iss": s.jwtTokenIssuer,
		"exp": &jwt.NumericDate{Time: time.Now().Add(time.Second * time.Duration(s.jwtTokenExpiresAfter))},
	})

	if token, err = tkn.SignedString(s.jwtSecretSalt); err != nil {
		return "", s.logger.LogPropagate(err)
	} else {
		return token, nil
	}
}

// Validate will decode the token and return a user ID or error, if it was occurred.
func (s *JwtService) Validate(token string) (userID vo.ID, err error) {
	// checking that token is not blocked
	found, err := s.blockedTokenRepository.Has(s.ctx, token)
	if err != nil {
		return vo.ID{}, s.logger.LogPropagate(err)
	}
	if found {
		return vo.ID{}, s.logger.LogPropagate(errors.NewAccessTokenWasBlockedError())
	}

	parsedToken, err := jwt.Parse(token, func(decodedToken *jwt.Token) (interface{}, error) {
		if decodedToken.Header["alg"] != s.jwtTokenEncryptAlgo {
			// user must be banned here because the algo wasn't matched
			return nil, errors.NewTokenAlgoWasNotMatchedError(token)
		}
		// cast to the configured givenToken signature type (stored in `s.jwtTokenEncryptAlgo`)
		if _, success := decodedToken.Method.(*jwt.SigningMethodHMAC); !success {
			return nil, errors.NewTokenUnexpectedSigningMethodError(token, decodedToken.Header["alg"])
		}
		// jwtSecretSalt is a string containing your secret, but you need pass the []byte
		return s.jwtSecretSalt, nil
	})
	if err != nil {
		// parsing givenToken error occurred
		return vo.ID{}, s.logger.LogPropagate(err)
	}

	// extracting claims of the givenToken payload
	if claims, success := parsedToken.Claims.(jwt.MapClaims); success && parsedToken.Valid {
		if err = s.isValidIssuer(token, claims); err != nil {
			return vo.ID{}, s.logger.LogPropagate(err)
		}

		userID, err = s.getUserID(claims)
		if err != nil {
			return vo.ID{}, s.logger.LogPropagate(err)
		}

		return userID, nil
	} else {
		// error occurred while extracting claims from givenToken or givenToken is not valid
		return vo.ID{}, errors.NewTokenInvalidError(token)
	}
}

// Block will mark the token as blocked into the storage.
func (s *JwtService) Block(token string) error {
	found, err := s.blockedTokenRepository.Has(s.ctx, token)
	if err != nil {
		return s.logger.LogPropagate(err)
	}
	if found {
		return s.logger.LogPropagate(errors.NewAccessTokenWasBlockedError())
	}
	return nil
}

func (s *JwtService) isValidIssuer(token string, claims jwt.Claims) error {
	// extracting the token issuer
	iss, err := claims.GetIssuer()
	if err != nil {
		return s.logger.LogPropagate(err)
	}
	// checking that token issuer is valid
	if iss != s.jwtTokenIssuer && !(func() (isIssuerWasMatched bool) {
		for _, acceptedIssuer := range s.jwtTokenAcceptedIssuers {
			if iss == acceptedIssuer {
				return true
			}
		}
		return false
	}()) {
		return errors.NewTokenIssuerWasNotMatchedError(token)
	}

	return nil
}

func (s *JwtService) getUserID(claims jwt.Claims) (userID vo.ID, err error) {
	// extracting subject (hexID) from the claims
	hexID, err := claims.GetSubject()
	if err != nil {
		return vo.ID{}, s.logger.LogPropagate(err)
	}
	// creating an object ID from hex
	oID, err := primitive.ObjectIDFromHex(hexID)
	if err != nil {
		return vo.ID{}, s.logger.LogPropagate(err)
	}
	// returning a success response
	return vo.ID{Value: oID}, nil
}
