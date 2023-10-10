package auth

import (
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

const AuthorizationPath = "/authorization"

// AuthorizationController - not implemented yet.
type AuthorizationController struct {
}

func NewAuthorizationController() *AuthorizationController {
	return &AuthorizationController{}
}

func (c *AuthorizationController) Create(w http.ResponseWriter, r *http.Request) {
	if _, err := w.Write([]byte("Sorry, the route is not implemented yet :(")); err != nil {
		log.Fatalln(err)
	}
}

func (c *AuthorizationController) AddRoute(router *mux.Router) {
	router.
		Path(AuthorizationPath).
		HandlerFunc(c.Create).
		Methods(http.MethodPost)
}
