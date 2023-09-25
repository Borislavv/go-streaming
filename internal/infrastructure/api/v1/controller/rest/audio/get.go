package audio

import (
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

const GetPath = "/audio/{id}"

type GetAudioController struct {
}

func NewGetController() *GetAudioController {
	return &GetAudioController{}
}

func (g *GetAudioController) Get(w http.ResponseWriter, r *http.Request) {
	if _, err := w.Write([]byte("Hello world from GET:item method!")); err != nil {
		log.Fatalln(err)
	}
}

func (g *GetAudioController) AddRoute(router *mux.Router) {
	router.
		Path(GetPath).
		HandlerFunc(g.Get).
		Methods(http.MethodGet)
}
