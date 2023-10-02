package user

import (
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

const DeletePath = "/user/{id}"

// DeleteUserController - not implemented yet.
type DeleteUserController struct {
}

func NewDeleteController() *DeleteUserController {
	return &DeleteUserController{}
}

func (d *DeleteUserController) Delete(w http.ResponseWriter, r *http.Request) {
	if _, err := w.Write([]byte("Sorry, the route is not implemented yet :(")); err != nil {
		log.Fatalln(err)
	}
}

func (d *DeleteUserController) AddRoute(router *mux.Router) {
	router.
		Path(DeletePath).
		HandlerFunc(d.Delete).
		Methods(http.MethodDelete)
}
