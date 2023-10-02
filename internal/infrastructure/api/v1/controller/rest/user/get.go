package user

import (
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

const GetPath = "/user/{id}"

// GetUserController - not implemented yet.
type GetUserController struct {
}

func NewGetController() *GetUserController {
	return &GetUserController{}
}

func (g *GetUserController) Get(w http.ResponseWriter, r *http.Request) {
	if _, err := w.Write([]byte("Sorry, the route is not implemented yet :(")); err != nil {
		log.Fatalln(err)
	}
}

func (g *GetUserController) AddRoute(router *mux.Router) {
	router.
		Path(GetPath).
		HandlerFunc(g.Get).
		Methods(http.MethodGet)
}
