package audio

import (
	"github.com/gorilla/mux"
	"net/http"
)

const CreatePath = "/audio"

// CreateAudioController - not implemented yet.
type CreateAudioController struct {
}

func NewCreateController() *CreateAudioController {
	return &CreateAudioController{}
}

func (c *CreateAudioController) Create(w http.ResponseWriter, r *http.Request) {
	if _, err := w.Write([]byte("Sorry, the route is not implemented yet :(")); err != nil {
		panic(err)
	}
}

func (c *CreateAudioController) AddRoute(router *mux.Router) {
	router.
		Path(CreatePath).
		HandlerFunc(c.Create).
		Methods(http.MethodPost)
}
