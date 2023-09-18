package render

import (
	"github.com/Borislavv/video-streaming/internal/infrastructure/api/v1/response"
	"github.com/Borislavv/video-streaming/internal/infrastructure/helper"
	"github.com/gorilla/mux"
	"html/template"
	"net/http"
)

const (
	IndexPath    = "/"
	TemplateName = "index.html"
)

type IndexController struct {
	responder response.Responder
}

func NewIndexController(responder response.Responder) *IndexController {
	return &IndexController{
		responder: responder,
	}
}

func (c *IndexController) Index(w http.ResponseWriter, r *http.Request) {
	tplPath, err := helper.TemplatePath(TemplateName)
	if err != nil {
		c.responder.Respond(w, err)
		return
	}

	tpl, err := template.ParseFiles(tplPath)
	if err != nil {
		c.responder.Respond(w, err)
		return
	}

	if err = tpl.Execute(w, nil); err != nil {
		if err != nil {
			c.responder.Respond(w, err)
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
