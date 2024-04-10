package auth

import (
	builderinterface "github.com/Borislavv/video-streaming/internal/domain/builder/interface"
	"github.com/Borislavv/video-streaming/internal/domain/logger/interface"
	authenticatorinterface "github.com/Borislavv/video-streaming/internal/domain/service/authenticator/interface"
	"github.com/Borislavv/video-streaming/internal/domain/service/di/interface"
	responseinterface "github.com/Borislavv/video-streaming/internal/infrastructure/api/v1/response/interface"
	"github.com/gorilla/mux"
	"net/http"
)

const AuthorizationPath = "/authorization"

type AuthorizationController struct {
	logger        loggerinterface.Logger
	builder       builderinterface.Auth
	authenticator authenticatorinterface.Authenticator
	responder     responseinterface.Responder
}

func NewAuthorizationController(serviceContainer diinterface.ServiceContainer) (*AuthorizationController, error) {
	loggerService, err := serviceContainer.GetLoggerService()
	if err != nil {
		return nil, err
	}

	authBuilder, err := serviceContainer.GetAuthBuilder()
	if err != nil {
		return nil, loggerService.LogPropagate(err)
	}

	authService, err := serviceContainer.GetAuthService()
	if err != nil {
		return nil, loggerService.LogPropagate(err)
	}

	responseService, err := serviceContainer.GetResponderService()
	if err != nil {
		return nil, loggerService.LogPropagate(err)
	}

	return &AuthorizationController{
		logger:        loggerService,
		builder:       authBuilder,
		authenticator: authService,
		responder:     responseService,
	}, nil
}

func (c *AuthorizationController) GetAccessToken(w http.ResponseWriter, r *http.Request) {
	// building an auth. request DTO
	req, err := c.builder.BuildAuthRequestDTOFromRequest(r)
	if err != nil {
		c.responder.Respond(w, c.logger.LogPropagate(err))
		return
	}

	// getting access token
	token, err := c.authenticator.Auth(req)
	if err != nil {
		c.responder.Respond(w, c.logger.LogPropagate(err))
		return
	}

	c.responder.Respond(w, token)
}

func (c *AuthorizationController) AddRoute(router *mux.Router) {
	router.
		Path(AuthorizationPath).
		HandlerFunc(c.GetAccessToken).
		Methods(http.MethodPost)
}
