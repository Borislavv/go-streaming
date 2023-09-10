package resource

import (
	"github.com/Borislavv/video-streaming/internal/domain/service"
	"github.com/Borislavv/video-streaming/internal/infrastructure/api/v1/response"
	"github.com/gorilla/mux"
	"net/http"
)

const UploadPath = "/resource"

type UploadResourceController struct {
	uploader service.Uploader
	responer response.Responder
}

func NewUploadResourceController(
	uploader service.Uploader,
	responer response.Responder,
) *UploadResourceController {
	return &UploadResourceController{
		uploader: uploader,
		responer: responer,
	}
}

func (c *UploadResourceController) Upload(w http.ResponseWriter, r *http.Request) {
	id, err := c.uploader.Upload(r)
	if err != nil {
		c.responer.Respond(w, err)
		return
	}
	c.responer.Respond(w, id)
}

func (c *UploadResourceController) AddRoute(router *mux.Router) {
	router.
		Path(UploadPath).
		HandlerFunc(c.Upload).
		Methods(http.MethodPost)
}
