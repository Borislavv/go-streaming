package auth

import (
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

const RegistrationPath = "/registration"

// RegistrationController - not implemented yet.
type RegistrationController struct {
}

func NewRegistrationController() *RegistrationController {
	return &RegistrationController{}
}

func (c *RegistrationController) Registration(w http.ResponseWriter, r *http.Request) {
	if _, err := w.Write([]byte("Sorry, the route is not implemented yet :(")); err != nil {
		log.Fatalln(err)
	}
}

func (c *RegistrationController) AddRoute(router *mux.Router) {
	router.
		Path(RegistrationPath).
		HandlerFunc(c.Registration).
		Methods(http.MethodPost)
}
