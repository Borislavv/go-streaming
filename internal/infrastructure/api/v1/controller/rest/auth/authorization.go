package auth

import (
	"github.com/Borislavv/video-streaming/internal/domain/builder"
	"github.com/Borislavv/video-streaming/internal/domain/logger"
	"github.com/Borislavv/video-streaming/internal/domain/service/auth"
	"github.com/Borislavv/video-streaming/internal/infrastructure/api/v1/response"
	"github.com/gorilla/mux"
	"net/http"
)

const AuthorizationPath = "/authorization"

type AuthorizationController struct {
	logger        logger.Logger
	builder       builder.Auth
	authenticator auth.Authenticator
	responder     response.Responder
}

func NewAuthorizationController(
	logger logger.Logger,
	builder builder.Auth,
	authenticator auth.Authenticator,
	responder response.Responder,
) *AuthorizationController {
	return &AuthorizationController{
		logger:        logger,
		builder:       builder,
		authenticator: authenticator,
		responder:     responder,
	}
}

func (c *AuthorizationController) Authorization(w http.ResponseWriter, r *http.Request) {
	reqDTO, err := c.builder.BuildAuthRequestDTOFromRequest(r)
	if err != nil {
		c.responder.Respond(w, c.logger.LogPropagate(err))
		return
	}

	token, err := c.authenticator.Auth(w, r, reqDTO)
	if err != nil {
		c.responder.Respond(w, c.logger.LogPropagate(err))
		return
	}

	c.responder.Respond(w, token)
}

func (c *AuthorizationController) AddRoute(router *mux.Router) {
	router.
		Path(AuthorizationPath).
		HandlerFunc(c.Authorization).
		Methods(http.MethodPost)
}
