package resource

import (
	"github.com/gorilla/mux"
	"net/http"
)

const UploadPath = "/upload"

type UploadResourceController struct {
}

func NewUploadResourceController() *UploadResourceController {
	return &UploadResourceController{}
}

func (c *UploadResourceController) Upload(w http.ResponseWriter, r *http.Request) {

}

func (c *UploadResourceController) AddRoute(router *mux.Router) {
	router.
		Path(UploadPath).
		HandlerFunc(c.Upload).
		Methods(http.MethodPost)
}
