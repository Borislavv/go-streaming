package resource

import (
	"github.com/Borislavv/video-streaming/internal/domain/builder"
	"github.com/Borislavv/video-streaming/internal/domain/logger"
	"github.com/Borislavv/video-streaming/internal/domain/service"
	"github.com/Borislavv/video-streaming/internal/infrastructure/api/v1/response"
	"github.com/gorilla/mux"
	"net/http"
)

const UploadPath = "/resource"

type UploadResourceController struct {
	logger    logger.Logger
	builder   builder.Resource
	service   service.Resource
	responder response.Responder
}

func NewUploadResourceController(
	logger logger.Logger,
	builder builder.Resource,
	service service.Resource,
	responer response.Responder,
) *UploadResourceController {
	return &UploadResourceController{
		logger:    logger,
		builder:   builder,
		service:   service,
		responder: responer,
	}
}

func (c *UploadResourceController) Upload(w http.ResponseWriter, r *http.Request) {
	req, err := c.builder.BuildUploadRequestDTOFromRequest(r)
	if err != nil {
		c.logger.Log(err)
		c.responder.Respond(w, err)
		return
	}

	agg, err := c.service.Upload(req)
	if err != nil {
		c.logger.Log(err)
		c.responder.Respond(w, err)
		return
	}

	c.responder.Respond(w, agg)
}

func (c *UploadResourceController) AddRoute(router *mux.Router) {
	router.
		Path(UploadPath).
		HandlerFunc(c.Upload).
		Methods(http.MethodPost)
}
