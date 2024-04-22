package render

import (
	"github.com/Borislavv/video-streaming/internal/domain/logger/interface"
	"github.com/Borislavv/video-streaming/internal/domain/service/di/interface"
	responseinterface "github.com/Borislavv/video-streaming/internal/infrastructure/api/v1/response/interface"
	"github.com/Borislavv/video-streaming/internal/infrastructure/helper"
	"github.com/gorilla/mux"
	"html/template"
	"net/http"
)

const (
	LoginPath        = "/login"
	authTemplateName = "auth.html"
)

type AuthController struct {
	logger    loggerinterface.Logger
	responder responseinterface.Responder
}

func NewAuthController(serviceContainer diinterface.ServiceContainer) (*AuthController, error) {
	loggerService, err := serviceContainer.GetLoggerService()
	if err != nil {
		return nil, err
	}

	responseService, err := serviceContainer.GetResponderService()
	if err != nil {
		return nil, loggerService.LogPropagate(err)
	}

	return &AuthController{
		logger:    loggerService,
		responder: responseService,
	}, nil
}

func (c *AuthController) Auth(w http.ResponseWriter, _ *http.Request) {
	tplPath, err := helper.TemplatePath(authTemplateName)
	if err != nil {
		c.responder.Respond(w, c.logger.LogPropagate(err))
		return
	}

	tpl, err := template.ParseFiles(tplPath)
	if err != nil {
		c.responder.Respond(w, c.logger.LogPropagate(err))
		return
	}

	if err = tpl.Execute(w, nil); err != nil {
		c.responder.Respond(w, c.logger.LogPropagate(err))
		return
	}
}

func (c *AuthController) AddRoute(router *mux.Router) {
	router.
		Path(LoginPath).
		HandlerFunc(c.Auth).
		Methods(http.MethodGet)
}
