package video

import (
	"github.com/Borislavv/video-streaming/internal/domain/builder"
	"github.com/Borislavv/video-streaming/internal/domain/service"
	"github.com/Borislavv/video-streaming/internal/infrastructure/api/v1/response"
	"github.com/gorilla/mux"
	"net/http"
)

const DeletePath = "/video/{id}"

type DeleteVideoController struct {
	builder  builder.Video
	service  service.Video
	response response.Responder
}

func NewDeleteVideoController(
	builder builder.Video,
	service service.Video,
	response response.Responder,
) *DeleteVideoController {
	return &DeleteVideoController{
		builder:  builder,
		service:  service,
		response: response,
	}
}

func (c *DeleteVideoController) Delete(w http.ResponseWriter, r *http.Request) {
	videoDto, err := c.builder.BuildDeleteRequestDTOFromRequest(r)
	if err != nil {
		c.response.Respond(w, err)
		return
	}

	if err = c.service.Delete(videoDto); err != nil {
		c.response.Respond(w, err)
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
