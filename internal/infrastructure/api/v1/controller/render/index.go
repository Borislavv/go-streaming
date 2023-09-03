package render

import (
	"github.com/Borislavv/video-streaming/internal/infrastructure/helper"
	"github.com/gorilla/mux"
	"html/template"
	"log"
	"net/http"
)

const (
	IndexPath    = "/"
	TemplateName = "index.html"
)

type IndexController struct {
}

func NewIndexController() *IndexController {
	return &IndexController{}
}

func (i *IndexController) Index(w http.ResponseWriter, r *http.Request) {
	tplPath, err := helper.TemplatePath(TemplateName)
	if err != nil {
		http.Error(w, "Internal server error, please contact with administrator.", http.StatusInternalServerError)
		log.Println("unable to parse path to index.html template", err)
		return
	}

	tpl, err := template.ParseFiles(tplPath)
	if err != nil {
		http.Error(w, "Internal server error, please contact with administrator.", http.StatusInternalServerError)
		log.Println("unable to parse index.html template", err)
		return
	}

	if err = tpl.Execute(w, nil); err != nil {
		if err != nil {
			http.Error(w, "Internal server error, please contact with administrator.", http.StatusInternalServerError)
			log.Println("unable to render index.html template", err)
			return
		}
	}
}

func (i *IndexController) AddRoute(router *mux.Router) {
	router.
		Path(IndexPath).
		HandlerFunc(i.Index).
		Methods(http.MethodGet)
}
