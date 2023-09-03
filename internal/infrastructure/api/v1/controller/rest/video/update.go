package video

import (
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

const UpdatePath = "/video/{id}"

type UpdateVideoController struct {
}

func NewUpdateVideoController() *UpdateVideoController {
	return &UpdateVideoController{}
}

func (u *UpdateVideoController) Update(w http.ResponseWriter, r *http.Request) {
	if _, err := w.Write([]byte("Hello world from PATCH method!")); err != nil {
		log.Fatalln(err)
	}
}

func (u *UpdateVideoController) AddRoute(router *mux.Router) {
	router.
		Path(UpdatePath).
		HandlerFunc(u.Update).
		Methods(http.MethodPatch)
}
