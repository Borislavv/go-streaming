package user

import (
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

const UpdatePath = "/audio/{id}"

// UpdateUserController - not implemented yet.
type UpdateUserController struct {
}

func NewUpdateUserController() *UpdateUserController {
	return &UpdateUserController{}
}

func (u *UpdateUserController) Update(w http.ResponseWriter, r *http.Request) {
	if _, err := w.Write([]byte("Sorry, the route is not implemented yet :(")); err != nil {
		log.Fatalln(err)
	}
}

func (u *UpdateUserController) AddRoute(router *mux.Router) {
	router.
		Path(UpdatePath).
		HandlerFunc(u.Update).
		Methods(http.MethodPatch)
}
