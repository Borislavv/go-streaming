package video

import (
	"github.com/Borislavv/video-streaming/internal/domain/builder"
	"github.com/Borislavv/video-streaming/internal/domain/service"
	"github.com/Borislavv/video-streaming/internal/infrastructure/api/v1/response"
	"github.com/gorilla/mux"
	"net/http"
)

const ListPath = "/video"

type ListVideoController struct {
	builder  builder.Video
	service  service.Video
	response response.Responder
}

func NewListVideoController(
	builder builder.Video,
	service service.Video,
	response response.Responder,
) *ListVideoController {
	return &ListVideoController{
		builder:  builder,
		service:  service,
		response: response,
	}
}

func (c *ListVideoController) List(w http.ResponseWriter, r *http.Request) {
	reqDto, err := c.builder.BuildListRequestDTOFromRequest(r)
	if err != nil {
		c.response.Respond(w, err)
		return
	}

	videos, err := c.service.List(reqDto)
	if err != nil {
		c.response.Respond(w, err)
		return
	}

	c.response.Respond(w, videos)
}

func (c *ListVideoController) AddRoute(router *mux.Router) {
	router.
		Path(ListPath).
		HandlerFunc(c.List).
		Methods(http.MethodGet)
}
