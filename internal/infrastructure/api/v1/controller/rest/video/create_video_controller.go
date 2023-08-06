package video

import (
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

const CreatePath = "/video"

type CreateVideoController struct {
}

func NewCreateController() *CreateVideoController {
	return &CreateVideoController{}
}

func (c *CreateVideoController) Create(w http.ResponseWriter, r *http.Request) {
	if _, err := w.Write([]byte("Hello world from POST method!")); err != nil {
		log.Fatalln(err)
	}
}

func (c *CreateVideoController) AddRoute(router *mux.Router) {
	router.
		Path(CreatePath).
		HandlerFunc(c.Create).
		Methods(http.MethodPost)
}
