package audio

import (
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

const ListPath = "/audio"

type ListAudioController struct {
}

func NewListVideoController() *ListAudioController {
	return &ListAudioController{}
}

func (l *ListAudioController) List(w http.ResponseWriter, r *http.Request) {
	if _, err := w.Write([]byte("Hello world from GET:list method!")); err != nil {
		log.Fatalln(err)
	}
}

func (l *ListAudioController) AddRoute(router *mux.Router) {
	router.
		Path(ListPath).
		HandlerFunc(l.List).
		Methods(http.MethodGet)
}
