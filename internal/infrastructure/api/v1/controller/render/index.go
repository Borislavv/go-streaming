package render

import (
	"github.com/Borislavv/video-streaming/internal/domain/logger/interface"
	"github.com/Borislavv/video-streaming/internal/domain/service/di/interface"
	response_interface "github.com/Borislavv/video-streaming/internal/infrastructure/api/v1/response/interface"
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
	logger    logger_interface.Logger
	responder response_interface.Responder
}

func NewIndexController(serviceContainer di_interface.ContainerManager) (*IndexController, error) {
	loggerService, err := serviceContainer.GetLoggerService()
	if err != nil {
		return nil, err
	}

	responseService, err := serviceContainer.GetResponderService()
	if err != nil {
		return nil, loggerService.LogPropagate(err)
	}

	return &IndexController{
		logger:    loggerService,
		responder: responseService,
	}, nil
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
