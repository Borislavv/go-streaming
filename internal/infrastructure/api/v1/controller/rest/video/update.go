package video

import (
	"github.com/Borislavv/video-streaming/internal/domain/builder"
	"github.com/Borislavv/video-streaming/internal/domain/logger"
	"github.com/Borislavv/video-streaming/internal/domain/service/video"
	"github.com/Borislavv/video-streaming/internal/infrastructure/api/v1/response"
	"github.com/gorilla/mux"
	"net/http"
)

const UpdatePath = "/video/{id}"

type UpdateVideoController struct {
	logger   logger.Logger
	builder  builder.Video
	service  video.CRUD
	response response.Responder
}

func NewUpdateVideoController(
	logger logger.Logger,
	builder builder.Video,
	service video.CRUD,
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
	videoDTO, err := c.builder.BuildUpdateRequestDTOFromRequest(r)
	if err != nil {
		c.response.Respond(w, c.logger.LogPropagate(err))
		return
	}

	videoAgg, err := c.service.Update(videoDTO)
	if err != nil {
		c.response.Respond(w, c.logger.LogPropagate(err))
		return
	}

	c.response.Respond(w, videoAgg)
}

func (c *UpdateVideoController) AddRoute(router *mux.Router) {
	router.
		Path(UpdatePath).
		HandlerFunc(c.Update).
		Methods(http.MethodPatch)
}
