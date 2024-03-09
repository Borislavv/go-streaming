package audio

import (
	"github.com/gorilla/mux"
	"net/http"
)

const ListPath = "/audio"

// ListAudioController - not implemented yet.
type ListAudioController struct {
}

func NewListController() *ListAudioController {
	return &ListAudioController{}
}

func (l *ListAudioController) List(w http.ResponseWriter, r *http.Request) {
	if _, err := w.Write([]byte("Sorry, the route is not implemented yet :(")); err != nil {
		panic(err)
	}
}

func (l *ListAudioController) AddRoute(router *mux.Router) {
	router.
		Path(ListPath).
		HandlerFunc(l.List).
		Methods(http.MethodGet)
}
