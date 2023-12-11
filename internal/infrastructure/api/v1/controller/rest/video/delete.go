package video

import (
	builder_interface "github.com/Borislavv/video-streaming/internal/domain/builder/interface"
	"github.com/Borislavv/video-streaming/internal/domain/logger/interface"
	di_interface "github.com/Borislavv/video-streaming/internal/domain/service/di/interface"
	video_interface "github.com/Borislavv/video-streaming/internal/domain/service/video/interface"
	response_interface "github.com/Borislavv/video-streaming/internal/infrastructure/api/v1/response/interface"
	"github.com/gorilla/mux"
	"net/http"
)

const DeletePath = "/video/{id}"

type DeleteController struct {
	logger    logger_interface.Logger
	builder   builder_interface.Video
	service   video_interface.CRUD
	responder response_interface.Responder
}

func NewDeleteController(serviceContainer di_interface.ContainerManager) (*DeleteController, error) {
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

	return &DeleteController{
		logger:    loggerService,
		builder:   videoBuilder,
		service:   videoCRUDService,
		responder: responseService,
	}, nil
}

func (c *DeleteController) Delete(w http.ResponseWriter, r *http.Request) {
	reqDTO, err := c.builder.BuildDeleteRequestDTOFromRequest(r)
	if err != nil {
		c.responder.Respond(w, c.logger.LogPropagate(err))
		return
	}

	if err = c.service.Delete(reqDTO); err != nil {
		c.responder.Respond(w, c.logger.LogPropagate(err))
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (c *DeleteController) AddRoute(router *mux.Router) {
	router.
		Path(DeletePath).
		HandlerFunc(c.Delete).
		Methods(http.MethodDelete)
}
