package video

import (
	"github.com/Borislavv/video-streaming/internal/domain/builder"
	"github.com/Borislavv/video-streaming/internal/domain/logger"
	"github.com/Borislavv/video-streaming/internal/domain/service/video"
	"github.com/Borislavv/video-streaming/internal/infrastructure/api/v1/response"
	"github.com/gorilla/mux"
	"net/http"
)

const ListPath = "/video"

type ListVideoController struct {
	logger   logger.Logger
	builder  builder.Video
	service  video.CRUD
	response response.Responder
}

func NewListController(
	logger logger.Logger,
	builder builder.Video,
	service video.CRUD,
	response response.Responder,
) *ListVideoController {
	return &ListVideoController{
		logger:   logger,
		builder:  builder,
		service:  service,
		response: response,
	}
}

func (c *ListVideoController) List(w http.ResponseWriter, r *http.Request) {
	reqDTO, e := c.builder.BuildListRequestDTOFromRequest(r)
	if e != nil {
		c.response.Respond(w, c.logger.LogPropagate(e))
		return
	}

	aggList, total, err := c.service.List(reqDTO)
	if err != nil {
		c.response.Respond(w, c.logger.LogPropagate(err))
		return
	}

	c.response.Respond(w,
		map[string]interface{}{
			"list": aggList,
			"pagination": map[string]interface{}{
				"page":  reqDTO.Page,
				"limit": reqDTO.Limit,
				"total": total,
			},
		},
	)
}

func (c *ListVideoController) AddRoute(router *mux.Router) {
	router.
		Path(ListPath).
		HandlerFunc(c.List).
		Methods(http.MethodGet)
}
