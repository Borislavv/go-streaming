package audio

import (
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

const DeletePath = "/audio/{id}"

type DeleteAudioController struct {
}

func NewDeleteVideoController() *DeleteAudioController {
	return &DeleteAudioController{}
}

func (d *DeleteAudioController) Delete(w http.ResponseWriter, r *http.Request) {
	if _, err := w.Write([]byte("Hello world from DELETE method!")); err != nil {
		log.Fatalln(err)
	}
}

func (d *DeleteAudioController) AddRoute(router *mux.Router) {
	router.
		Path(DeletePath).
		HandlerFunc(d.Delete).
		Methods(http.MethodDelete)
}
