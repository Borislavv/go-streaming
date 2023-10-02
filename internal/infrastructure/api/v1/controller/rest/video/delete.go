package video

import (
	"github.com/Borislavv/video-streaming/internal/domain/builder"
	"github.com/Borislavv/video-streaming/internal/domain/logger"
	"github.com/Borislavv/video-streaming/internal/domain/service/video"
	"github.com/Borislavv/video-streaming/internal/infrastructure/api/v1/response"
	"github.com/gorilla/mux"
	"net/http"
)

const DeletePath = "/video/{id}"

type DeleteVideoController struct {
	logger   logger.Logger
	builder  builder.Video
	service  video.CRUD
	response response.Responder
}

func NewDeleteController(
	logger logger.Logger,
	builder builder.Video,
	service video.CRUD,
	response response.Responder,
) *DeleteVideoController {
	return &DeleteVideoController{
		logger:   logger,
		builder:  builder,
		service:  service,
		response: response,
	}
}

func (c *DeleteVideoController) Delete(w http.ResponseWriter, r *http.Request) {
	reqDTO, err := c.builder.BuildDeleteRequestDTOFromRequest(r)
	if err != nil {
		c.response.Respond(w, c.logger.LogPropagate(err))
		return
	}

	if err = c.service.Delete(reqDTO); err != nil {
		c.response.Respond(w, c.logger.LogPropagate(err))
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (c *DeleteVideoController) AddRoute(router *mux.Router) {
	router.
		Path(DeletePath).
		HandlerFunc(c.Delete).
		Methods(http.MethodDelete)
}
