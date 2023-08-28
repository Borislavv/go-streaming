package video

import (
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

const ListPath = "/video"

type ListVideoController struct {
}

func NewListVideoController() *ListVideoController {
	return &ListVideoController{}
}

func (l *ListVideoController) List(w http.ResponseWriter, r *http.Request) {
	if _, err := w.Write([]byte("Hello world from GET:list method!")); err != nil {
		log.Fatalln(err)
	}
}

func (l *ListVideoController) AddRoute(router *mux.Router) {
	router.
		Path(ListPath).
		HandlerFunc(l.List).
		Methods(http.MethodGet)
}
