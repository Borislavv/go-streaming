package resource

import (
	"github.com/Borislavv/video-streaming/internal/domain/builder"
	"github.com/Borislavv/video-streaming/internal/domain/logger"
	"github.com/Borislavv/video-streaming/internal/domain/service"
	"github.com/Borislavv/video-streaming/internal/infrastructure/api/v1/response"
	"net/http"
)

type ListResourceController struct {
	logger    logger.Logger
	builder   builder.Resource
	service   service.Resource
	responder response.Responder
}

func NewListResourceController(
	logger logger.Logger,
	builder builder.Resource,
	service service.Resource,
	responder response.Responder,
) *ListResourceController {
	return &ListResourceController{
		logger:    logger,
		builder:   builder,
		service:   service,
		responder: responder,
	}
}

func (c *ListResourceController) List(w http.ResponseWriter, r *http.Request) {
	
}
