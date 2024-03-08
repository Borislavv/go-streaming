package audio

import (
	"github.com/gorilla/mux"
	"net/http"
)

const GetPath = "/audio/{id}"

// GetAudioController - not implemented yet.
type GetAudioController struct {
}

func NewGetController() *GetAudioController {
	return &GetAudioController{}
}

func (g *GetAudioController) Get(w http.ResponseWriter, r *http.Request) {
	if _, err := w.Write([]byte("Sorry, the route is not implemented yet :(")); err != nil {
		panic(err)
	}
}

func (g *GetAudioController) AddRoute(router *mux.Router) {
	router.
		Path(GetPath).
		HandlerFunc(g.Get).
		Methods(http.MethodGet)
}
