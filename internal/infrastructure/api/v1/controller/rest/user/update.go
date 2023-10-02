package user

import (
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

const UpdatePath = "/audio/{id}"

// UpdateController - not implemented yet.
type UpdateController struct {
}

func NewUpdateUserController() *UpdateController {
	return &UpdateController{}
}

func (u *UpdateController) Update(w http.ResponseWriter, r *http.Request) {
	if _, err := w.Write([]byte("Sorry, the route is not implemented yet :(")); err != nil {
		log.Fatalln(err)
	}
}

func (u *UpdateController) AddRoute(router *mux.Router) {
	router.
		Path(UpdatePath).
		HandlerFunc(u.Update).
		Methods(http.MethodPatch)
}
