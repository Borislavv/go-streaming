package video

import (
	"github.com/Borislavv/video-streaming/internal/domain/builder"
	"github.com/Borislavv/video-streaming/internal/domain/logger"
	"github.com/Borislavv/video-streaming/internal/domain/service"
	"github.com/Borislavv/video-streaming/internal/infrastructure/api/v1/response"
	"github.com/gorilla/mux"
	"net/http"
)

const UpdatePath = "/video/{id}"

type UpdateVideoController struct {
	logger   logger.Logger
	builder  builder.Video
	service  service.Video
	response response.Responder
}

func NewUpdateVideoController(
	logger logger.Logger,
	builder builder.Video,
	service service.Video,
	response response.Responder,
) *UpdateVideoController {
	return &UpdateVideoController{
		logger:   logger,
		builder:  builder,
		service:  service,
		response: response,
	}
}

func (c *UpdateVideoController) Update(w http.ResponseWriter, r *http.Request) {
	videoDto, err := c.builder.BuildUpdateRequestDTOFromRequest(r)
	if err != nil {
		c.response.Respond(w, c.logger.LogPropagate(err))
		return
	}

	video, err := c.service.Update(videoDto)
	if err != nil {
		c.response.Respond(w, c.logger.LogPropagate(err))
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
