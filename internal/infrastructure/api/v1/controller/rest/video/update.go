package video

import (
	"github.com/Borislavv/video-streaming/internal/domain/builder/interface"
	"github.com/Borislavv/video-streaming/internal/domain/logger/interface"
	"github.com/Borislavv/video-streaming/internal/domain/service/di/interface"
	video_interface "github.com/Borislavv/video-streaming/internal/domain/service/video/interface"
	response_interface "github.com/Borislavv/video-streaming/internal/infrastructure/api/v1/response/interface"
	"github.com/gorilla/mux"
	"net/http"
)

const UpdatePath = "/video/{id}"

type UpdateController struct {
	logger   logger_interface.Logger
	builder  builder_interface.Video
	service  video_interface.CRUD
	response response_interface.Responder
}

func NewUpdateController(serviceContainer diinterface.ContainerManager) (*UpdateController, error) {
	loggerService, err := serviceContainer.GetLoggerService()
	if err != nil {
		return nil, err
	}

	videoBuilder, err := serviceContainer.GetVideoBuilder()
	if err != nil {
		return nil, loggerService.LogPropagate(err)
	}

	videoCRUDService, err := serviceContainer.GetVideoCRUDService()
	if err != nil {
		return nil, loggerService.LogPropagate(err)
	}

	responseService, err := serviceContainer.GetResponderService()
	if err != nil {
		return nil, loggerService.LogPropagate(err)
	}

	return &UpdateController{
		logger:   loggerService,
		builder:  videoBuilder,
		service:  videoCRUDService,
		response: responseService,
	}, nil
}

func (c *UpdateController) Update(w http.ResponseWriter, r *http.Request) {
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

func (c *UpdateController) AddRoute(router *mux.Router) {
	router.
		Path(UpdatePath).
		HandlerFunc(c.Update).
		Methods(http.MethodPatch)
}
