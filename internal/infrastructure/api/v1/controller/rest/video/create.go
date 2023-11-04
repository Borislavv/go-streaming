package video

import (
	"github.com/Borislavv/video-streaming/internal/domain/builder"
	"github.com/Borislavv/video-streaming/internal/domain/logger"
	"github.com/Borislavv/video-streaming/internal/domain/service/authenticator"
	"github.com/Borislavv/video-streaming/internal/domain/service/video"
	"github.com/Borislavv/video-streaming/internal/infrastructure/api/v1/response"
	"github.com/gorilla/mux"
	"net/http"
)

const CreatePath = "/video"

type CreateVideoController struct {
	logger      logger.Logger
	builder     builder.Video
	crudService video.CRUD
	authService authenticator.Authenticator
	response    response.Responder
}

func NewCreateController(
	logger logger.Logger,
	builder builder.Video,
	crudService video.CRUD,
	authService authenticator.Authenticator,
	response response.Responder,
) *CreateVideoController {
	return &CreateVideoController{
		logger:      logger,
		builder:     builder,
		crudService: crudService,
		authService: authService,
		response:    response,
	}
}

func (c *CreateVideoController) Create(w http.ResponseWriter, r *http.Request) {
	userID, err := c.authService.IsAuthed(r)
	if err != nil {
		c.response.Respond(w, c.logger.LogPropagate(err))
		return
	}

	videoDTO, err := c.builder.BuildCreateRequestDTOFromRequest(r)
	if err != nil {
		c.response.Respond(w, c.logger.LogPropagate(err))
		return
	}

	videoAgg, err := c.crudService.Create(videoDTO)
	if err != nil {
		c.response.Respond(w, c.logger.LogPropagate(err))
		return
	}

	c.response.Respond(w, videoAgg)
	w.WriteHeader(http.StatusCreated)
}

func (c *CreateVideoController) AddRoute(router *mux.Router) {
	router.
		Path(CreatePath).
		HandlerFunc(c.Create).
		Methods(http.MethodPost)
}
