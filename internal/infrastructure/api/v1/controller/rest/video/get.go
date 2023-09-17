package video

import (
	"github.com/Borislavv/video-streaming/internal/domain/builder"
	"github.com/Borislavv/video-streaming/internal/domain/service"
	"github.com/Borislavv/video-streaming/internal/infrastructure/api/v1/response"
	"github.com/gorilla/mux"
	"net/http"
)

const GetPath = "/video/{id}"

type GetVideoController struct {
	builder  builder.Video
	service  service.Video
	response response.Responder
}

func NewGetVideoController(
	builder builder.Video,
	service service.Video,
	response response.Responder,
) *GetVideoController {
	return &GetVideoController{
		builder:  builder,
		service:  service,
		response: response,
	}
}

func (c *GetVideoController) Get(w http.ResponseWriter, r *http.Request) {
	req, err := c.builder.BuildGetRequestDTOFromRequest(r)
	if err != nil {
		c.response.Respond(w, err)
		return
	}

	video, err := c.service.Get(req)
	if err != nil {
		c.response.Respond(w, err)
		return
	}

	c.response.Respond(w, video)
}

func (c *GetVideoController) AddRoute(router *mux.Router) {
	router.
		Path(GetPath).
		HandlerFunc(c.Get).
		Methods(http.MethodGet)
}
