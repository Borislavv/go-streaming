package resource

import (
	"github.com/Borislavv/video-streaming/internal/domain/builder"
	"github.com/Borislavv/video-streaming/internal/domain/logger"
	"github.com/Borislavv/video-streaming/internal/domain/service/resource"
	"github.com/Borislavv/video-streaming/internal/infrastructure/api/v1/response"
	"github.com/gorilla/mux"
	"net/http"
)

const UploadPath = "/resource"

type UploadResourceController struct {
	logger    logger.Logger
	builder   builder.Resource
	service   resource.CRUD
	responder response.Responder
}

func NewUploadController(
	logger logger.Logger,
	builder builder.Resource,
	service resource.CRUD,
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
	reqDTO, err := c.builder.BuildUploadRequestDTOFromRequest(r)
	if err != nil {
		c.responder.Respond(w, c.logger.LogPropagate(err))
		return
	}

	resourceAgg, err := c.service.Upload(reqDTO)
	if err != nil {
		c.responder.Respond(w, c.logger.LogPropagate(err))
		return
	}

	c.responder.Respond(w, resourceAgg)
}

func (c *UploadResourceController) AddRoute(router *mux.Router) {
	router.
		Path(UploadPath).
		HandlerFunc(c.Upload).
		Methods(http.MethodPost)
}
