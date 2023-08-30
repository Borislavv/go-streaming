package video

import (
	"github.com/Borislavv/video-streaming/internal/app/service"
	"github.com/Borislavv/video-streaming/internal/app/service/logger"
	dtobuilder "github.com/Borislavv/video-streaming/internal/domain/builder/dto"
	"github.com/gorilla/mux"
	"net/http"
)

const CreatePath = "/video"

type CreateVideoController struct {
	logger          logger.Logger
	videoDtoBuilder dtobuilder.Video
	videoCreator    service.VideoCreator
}

func NewCreateController(
	logger logger.Logger,
	videoDtoBuilder dtobuilder.Video,
	videoCreator service.VideoCreator,
) *CreateVideoController {
	return &CreateVideoController{
		logger:          logger,
		videoDtoBuilder: videoDtoBuilder,
		videoCreator:    videoCreator,
	}
}

func (c *CreateVideoController) Create(w http.ResponseWriter, r *http.Request) {
	videoDto, err := c.videoDtoBuilder.BuildFromRequest(r)
	if err != nil {
		c.logger.Error(err)
		// return the error response
		if _, err = w.Write([]byte("Internal server error, please contact with service administrator.")); err != nil {
			c.logger.Critical(err)
		}
		return
	}

	id, err := c.videoCreator.Create(videoDto)
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
