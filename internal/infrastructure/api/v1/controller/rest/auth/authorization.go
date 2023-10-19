package auth

import (
	"fmt"
	"github.com/Borislavv/video-streaming/internal/domain/logger"
	"github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"time"
)

const AuthorizationPath = "/authorization"

type AuthorizationController struct {
	logger logger.Logger
}

func NewAuthorizationController(logger logger.Logger) *AuthorizationController {
	return &AuthorizationController{logger: logger}
}

var tokenStr string
var hmacSampleSecret []byte = []byte("jared-streaming-service")

func (c *AuthorizationController) Authorization(w http.ResponseWriter, r *http.Request) {
	rcookie, err := r.Cookie("access-token")
	if err != nil {
		c.logger.Error("access-token cookie is not present into request")
	} else {
		if rcookie.Value == tokenStr {
			token, err := jwt.Parse(rcookie.Value, func(token *jwt.Token) (interface{}, error) {
				// Don't forget to validate the alg is what you expect:
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, c.logger.ErrorPropagate(fmt.Sprintf("Unexpected signing method: %v", token.Header["alg"]))
				}

				// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
				return hmacSampleSecret, nil
			})
			if err != nil {
				c.logger.Error(err)
				return
			}
			if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
				if _, err = w.Write([]byte(fmt.Sprintf("sub: %v, iss: %v, exp: %v", claims["sub"], claims["iss"], claims["exp"]))); err != nil {
					c.logger.Error(err)
					return
				}
			} else {
				c.logger.Error(err)
				return
			}
		} else {
			c.logger.Error("tokenStr != rcookie.Value error")
			return
		}
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": "jared-user",
		"iss": "streaming-service",
		"exp": time.Now().Add(time.Minute * 2),
	})

	tokenStr, err = token.SignedString(hmacSampleSecret)
	if err != nil {
		log.Fatalln(err)
	}

	cookie := &http.Cookie{
		Name:     "access-token",
		Value:    tokenStr,
		HttpOnly: true,
	}

	http.SetCookie(w, cookie)

	if _, err = w.Write([]byte("cookie successfully sat up")); err != nil {
		log.Fatalln(err)
	}
}

func (c *AuthorizationController) AddRoute(router *mux.Router) {
	router.
		Path(AuthorizationPath).
		HandlerFunc(c.Authorization).
		Methods(http.MethodGet)
}
