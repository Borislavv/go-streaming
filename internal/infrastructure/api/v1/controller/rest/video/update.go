package video

import (
	"github.com/Borislavv/video-streaming/internal/domain/builder"
	"github.com/Borislavv/video-streaming/internal/domain/service"
	"github.com/Borislavv/video-streaming/internal/infrastructure/api/v1/response"
	"github.com/gorilla/mux"
	"net/http"
)

const UpdatePath = "/video/{id}"

type UpdateVideoController struct {
	builder  builder.Video
	service  service.Video
	response response.Responder
}

func NewUpdateVideoController(
	builder builder.Video,
	service service.Video,
	response response.Responder,
) *UpdateVideoController {
	return &UpdateVideoController{
		builder:  builder,
		service:  service,
		response: response,
	}
}

func (c *UpdateVideoController) Update(w http.ResponseWriter, r *http.Request) {
	videoDto, err := c.builder.BuildUpdateRequestDtoFromRequest(r)
	if err != nil {
		c.response.Respond(w, err)
		return
	}

	video, err := c.service.Update(videoDto)
	if err != nil {
		c.response.Respond(w, err)
		return
	}

	c.response.Respond(w, video)
}

func (c *UpdateVideoController) AddRoute(router *mux.Router) {
	router.
		Path(UpdatePath).
		HandlerFunc(c.Update).
		Methods(http.MethodPatch)
}
