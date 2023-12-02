package render

import (
	"github.com/Borislavv/video-streaming/internal/domain/logger"
	"github.com/Borislavv/video-streaming/internal/infrastructure/api/v1/response"
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
	logger    logger.Logger
	responder response.Responder
}

func NewLoginController(
	logger logger.Logger,
	responder response.Responder,
) *LoginController {
	return &LoginController{
		logger:    logger,
		responder: responder,
	}
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
