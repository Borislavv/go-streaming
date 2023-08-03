package video

import (
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

const DeletePath = "/video/{id}"

type DeleteVideoController struct {
}

func NewDeleteVideoController() *DeleteVideoController {
	return &DeleteVideoController{}
}

func (c *DeleteVideoController) Delete(w http.ResponseWriter, r *http.Request) {
	if _, err := w.Write([]byte("Hello world from DELETE method!")); err != nil {
		log.Fatalln(err)
	}
}

func (c *DeleteVideoController) AddRoute(router *mux.Router) {
	router.
		Path(DeletePath).
		HandlerFunc(c.Delete).
		Methods(http.MethodDelete)
}
