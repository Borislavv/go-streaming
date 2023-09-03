package video

import (
	"github.com/Borislavv/video-streaming/internal/domain/builder"
	"github.com/Borislavv/video-streaming/internal/domain/service"
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
	videoDto, err := c.builder.BuildCreateRequestDtoFromRequest(r)
	if err != nil {
		c.logger.Error(err)
		// return the error response
		if _, err = w.Write([]byte("Internal server error, please contact with service administrator.")); err != nil {
			c.logger.Critical(err)
		}
		return
	}

	id, err := c.service.Create(videoDto)
	if err != nil {
		c.logger.Error(err)
		// return the error response
		if _, err = w.Write([]byte("Internal server error, please contact with service administrator.")); err != nil {
			c.logger.Critical(err)
		}
		return
	}

	if _, err = w.Write([]byte(id)); err != nil {
		c.logger.Critical(err)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

func (c *CreateVideoController) AddRoute(router *mux.Router) {
	router.
		Path(CreatePath).
		HandlerFunc(c.Create).
		Methods(http.MethodPost)
}
