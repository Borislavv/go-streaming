package audio

import (
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

const DeletePath = "/audio/{id}"

// DeleteAudioController - not implemented yet.
type DeleteAudioController struct {
}

func NewDeleteController() *DeleteAudioController {
	return &DeleteAudioController{}
}

func (d *DeleteAudioController) Delete(w http.ResponseWriter, r *http.Request) {
	if _, err := w.Write([]byte("Sorry, the route is not implemented yet :(")); err != nil {
		log.Fatalln(err)
	}
}

func (d *DeleteAudioController) AddRoute(router *mux.Router) {
	router.
		Path(DeletePath).
		HandlerFunc(d.Delete).
		Methods(http.MethodDelete)
}
