package audio

import (
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

const CreatePath = "/audio"

type CreateAudioController struct {
}

func NewCreateController() *CreateAudioController {
	return &CreateAudioController{}
}

func (c *CreateAudioController) Create(w http.ResponseWriter, r *http.Request) {
	if _, err := w.Write([]byte("Hello world from POST method!")); err != nil {
		log.Fatalln(err)
	}
}

func (c *CreateAudioController) AddRoute(router *mux.Router) {
	router.
		Path(CreatePath).
		HandlerFunc(c.Create).
		Methods(http.MethodPost)
}
