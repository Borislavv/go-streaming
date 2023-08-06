package static

import (
	"github.com/Borislavv/video-streaming/internal/infrastructure/helper"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"strings"
)

const ResourcesPrefix = "/static/"

type ResourceController struct {
}

func NewResourceController() *ResourceController {
	return &ResourceController{}
}

func (i *ResourceController) Serve(w http.ResponseWriter, r *http.Request) {
	dir, err := helper.ResourcesDir()
	if err != nil {
		http.Error(w, "Internal server error, please contact with administrator.", http.StatusInternalServerError)
		log.Println("unable to serve static files due to unable receive resources path", err)
		return
	}

	path := r.URL.Path
	if strings.Contains(path, ResourcesPrefix) {
		path = strings.ReplaceAll(path, ResourcesPrefix, "")
	}

	http.ServeFile(w, r, dir+path)
}

func (i *ResourceController) AddRoute(router *mux.Router) {
	router.
		PathPrefix(ResourcesPrefix).
		HandlerFunc(i.Serve).
		Methods(http.MethodGet)
}
