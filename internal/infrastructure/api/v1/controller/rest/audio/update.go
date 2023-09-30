package audio

import (
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

const UpdatePath = "/audio/{id}"

// UpdateAudioController - not implemented yet.
type UpdateAudioController struct {
}

func NewUpdateController() *UpdateAudioController {
	return &UpdateAudioController{}
}

func (u *UpdateAudioController) Update(w http.ResponseWriter, r *http.Request) {
	if _, err := w.Write([]byte("Hello world from PATCH method!")); err != nil {
		log.Fatalln(err)
	}
}

func (u *UpdateAudioController) AddRoute(router *mux.Router) {
	router.
		Path(UpdatePath).
		HandlerFunc(u.Update).
		Methods(http.MethodPatch)
}
