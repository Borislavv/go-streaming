package video

import (
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

const GetPath = "/video/{id}"

type GetVideoController struct {
}

func NewGetVideoController() *GetVideoController {
	return &GetVideoController{}
}

func (g *GetVideoController) Get(w http.ResponseWriter, r *http.Request) {
	if _, err := w.Write([]byte("Hello world from GET:item method!")); err != nil {
		log.Fatalln(err)
	}
}

func (g *GetVideoController) AddRoute(router *mux.Router) {
	router.
		Path(GetPath).
		HandlerFunc(g.Get).
		Methods(http.MethodGet)
}
