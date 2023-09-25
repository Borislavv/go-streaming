package video

import (
	"github.com/Borislavv/video-streaming/internal/domain/builder"
	"github.com/Borislavv/video-streaming/internal/domain/logger"
	"github.com/Borislavv/video-streaming/internal/domain/service/video"
	"github.com/Borislavv/video-streaming/internal/infrastructure/api/v1/response"
	"github.com/gorilla/mux"
	"net/http"
)

const GetPath = "/video/{id}"

type GetVideoController struct {
	logger   logger.Logger
	builder  builder.Video
	service  video.CRUD
	response response.Responder
}

func NewGetController(
	logger logger.Logger,
	builder builder.Video,
	service video.CRUD,
	response response.Responder,
) *GetVideoController {
	return &GetVideoController{
		logger:   logger,
		builder:  builder,
		service:  service,
		response: response,
	}
}

func (c *GetVideoController) Get(w http.ResponseWriter, r *http.Request) {
	reqDTO, err := c.builder.BuildGetRequestDTOFromRequest(r)
	if err != nil {
		c.response.Respond(w, c.logger.LogPropagate(err))
		return
	}

	videoAgg, err := c.service.Get(reqDTO)
	if err != nil {
		c.response.Respond(w, c.logger.LogPropagate(err))
		return
	}

	c.response.Respond(w, videoAgg)
}

func (c *GetVideoController) AddRoute(router *mux.Router) {
	router.
		Path(GetPath).
		HandlerFunc(c.Get).
		Methods(http.MethodGet)
}
