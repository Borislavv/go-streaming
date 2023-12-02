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
	IndexPath         = "/"
	IndexTemplateName = "index.html"
)

type IndexController struct {
	logger    logger.Logger
	responder response.Responder
}

func NewIndexController(
	logger logger.Logger,
	responder response.Responder,
) *IndexController {
	return &IndexController{
		logger:    logger,
		responder: responder,
	}
}

func (c *IndexController) Index(w http.ResponseWriter, _ *http.Request) {
	tplPath, err := helper.TemplatePath(IndexTemplateName)
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

func (c *IndexController) AddRoute(router *mux.Router) {
	router.
		Path(IndexPath).
		HandlerFunc(c.Index).
		Methods(http.MethodGet)
}
