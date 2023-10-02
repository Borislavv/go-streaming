package user

import (
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

const CreatePath = "/user"

// CreateUserController - not implemented yet.
type CreateUserController struct {
}

func NewCreateController() *CreateUserController {
	return &CreateUserController{}
}

func (c *CreateUserController) Create(w http.ResponseWriter, r *http.Request) {
	if _, err := w.Write([]byte("Sorry, the route is not implemented yet :(")); err != nil {
		log.Fatalln(err)
	}
}

func (c *CreateUserController) AddRoute(router *mux.Router) {
	router.
		Path(CreatePath).
		HandlerFunc(c.Create).
		Methods(http.MethodPost)
}
