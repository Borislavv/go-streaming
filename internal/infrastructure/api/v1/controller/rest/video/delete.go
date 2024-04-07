package video

import (
	builderinterface "github.com/Borislavv/video-streaming/internal/domain/builder/interface"
	"github.com/Borislavv/video-streaming/internal/domain/logger/interface"
	diinterface "github.com/Borislavv/video-streaming/internal/domain/service/di/interface"
	videointerface "github.com/Borislavv/video-streaming/internal/domain/service/video/interface"
	responseinterface "github.com/Borislavv/video-streaming/internal/infrastructure/api/v1/response/interface"
	"github.com/gorilla/mux"
	"net/http"
)

const DeletePath = "/video/{id}"

type DeleteController struct {
	logger    loggerinterface.Logger
	builder   builderinterface.Video
	service   videointerface.CRUD
	responder responseinterface.Responder
}

func NewDeleteController(serviceContainer diinterface.ServiceContainer) (*DeleteController, error) {
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
