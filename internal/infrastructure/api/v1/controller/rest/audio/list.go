package audio

import (
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

const ListPath = "/audio"

// ListAudioController - not implemented yet.
type ListAudioController struct {
}

func NewListAudioController() *ListAudioController {
	return &ListAudioController{}
}

func (l *ListAudioController) List(w http.ResponseWriter, r *http.Request) {
	if _, err := w.Write([]byte("Sorry, the route is not implemented yet :(")); err != nil {
		log.Fatalln(err)
	}
}

func (l *ListAudioController) AddRoute(router *mux.Router) {
	router.
		Path(ListPath).
		HandlerFunc(l.List).
		Methods(http.MethodGet)
}
