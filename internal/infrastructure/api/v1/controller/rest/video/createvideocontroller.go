package video

import (
	"context"
	"github.com/Borislavv/video-streaming/internal/app/service/logger"
	aggbuilder "github.com/Borislavv/video-streaming/internal/domain/builder/agg"
	dtobuilder "github.com/Borislavv/video-streaming/internal/domain/builder/dto"
	entitybuilder "github.com/Borislavv/video-streaming/internal/domain/builder/entity"
	"github.com/Borislavv/video-streaming/internal/domain/repository"
	"github.com/gorilla/mux"
	"net/http"
)

const CreatePath = "/video"

type CreateVideoController struct {
	logger             logger.Logger
	videoDtoBuilder    dtobuilder.Video
	videoEntityBuilder entitybuilder.Video
	videoAggBuilder    aggbuilder.Video
	videoRepository    repository.VideoRepository
}

func NewCreateController(
	logger logger.Logger,
	videoDtoBuilder dtobuilder.Video,
	videoEntityBuilder entitybuilder.Video,
	videoAggBuilder aggbuilder.Video,
	videoRepository repository.VideoRepository,
) *CreateVideoController {
	return &CreateVideoController{
		logger:             logger,
		videoDtoBuilder:    videoDtoBuilder,
		videoEntityBuilder: videoEntityBuilder,
		videoAggBuilder:    videoAggBuilder,
		videoRepository:    videoRepository,
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

	id, err := c.videoRepository.Insert(
		context.Background(),
		c.videoAggBuilder.Build(
			c.videoEntityBuilder.Build(videoDto),
		),
	)
	if err != nil {
		c.logger.Error(err)
		// return the error response
		if _, err = w.Write([]byte("Internal server error, please contact with service administrator.")); err != nil {
			c.logger.Critical(err)
		}
		return
	}

	if _, err = w.Write([]byte("Created. Id: " + id)); err != nil {
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
