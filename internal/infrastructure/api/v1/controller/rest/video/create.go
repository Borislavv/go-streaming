package video

import (
	"github.com/Borislavv/video-streaming/internal/domain/builder"
	"github.com/Borislavv/video-streaming/internal/domain/service"
	"github.com/Borislavv/video-streaming/internal/infrastructure/api/v1/response"
	"github.com/gorilla/mux"
	"net/http"
)

const CreatePath = "/video"

type CreateVideoController struct {
	logger  service.Logger
	builder builder.Video
	service service.Video
}

func NewCreateController(
	logger service.Logger,
	builder builder.Video,
	service service.Video,
) *CreateVideoController {
	return &CreateVideoController{
		logger:  logger,
		builder: builder,
		service: service,
	}
}

func (c *CreateVideoController) Create(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	videoDto, err := c.builder.BuildCreateRequestDtoFromRequest(r)
	if err != nil {
		c.logger.Error(err)
		// return the error response
		bytes, rerr := response.Respond(nil, err)
		if err != nil {
			c.logger.Critical(rerr)
		}
		if _, err = w.Write(bytes); err != nil {
			c.logger.Critical(err)
		}
		return
	}

	id, err := c.service.Create(videoDto)
	if err != nil {
		c.logger.Error(err)
		// return the error response
		bytes, rerr := response.Respond(nil, err)
		if rerr != nil {
			c.logger.Critical(rerr)
			return
		}
		if _, err = w.Write(bytes); err != nil {
			c.logger.Critical(err)
		}
		return
	}

	bytes, err := response.Respond(id, nil)
	if err != nil {
		c.logger.Critical(err)
		return
	}
	w.WriteHeader(http.StatusCreated)
	if _, err = w.Write(bytes); err != nil {
		c.logger.Critical(err)
		return
	}
}

func (c *CreateVideoController) AddRoute(router *mux.Router) {
	router.
		Path(CreatePath).
		HandlerFunc(c.Create).
		Methods(http.MethodPost)
}
