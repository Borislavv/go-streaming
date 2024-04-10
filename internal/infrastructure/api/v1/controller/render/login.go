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
	LoginPath         = "/login"
	LoginTemplateName = "login.html"
)

type LoginController struct {
	logger    loggerinterface.Logger
	responder responseinterface.Responder
}

func NewLoginController(serviceContainer diinterface.ServiceContainer) (*LoginController, error) {
	loggerService, err := serviceContainer.GetLoggerService()
	if err != nil {
		return nil, err
	}

	responseService, err := serviceContainer.GetResponderService()
	if err != nil {
		return nil, loggerService.LogPropagate(err)
	}

	return &LoginController{
		logger:    loggerService,
		responder: responseService,
	}, nil
}

func (c *LoginController) Login(w http.ResponseWriter, _ *http.Request) {
	tplPath, err := helper.TemplatePath(LoginTemplateName)
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
		if err != nil {
			c.responder.Respond(w, c.logger.LogPropagate(err))
			return
		}
	}
}

func (c *LoginController) AddRoute(router *mux.Router) {
	router.
		Path(LoginPath).
		HandlerFunc(c.Login).
		Methods(http.MethodGet)
}
