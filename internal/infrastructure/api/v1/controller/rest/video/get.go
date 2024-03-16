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

const GetPath = "/video/{id}"

type GetController struct {
	logger    loggerinterface.Logger
	builder   builder_interface.Video
	service   video_interface.CRUD
	responder response_interface.Responder
}

func NewGetController(serviceContainer diinterface.ContainerManager) (*GetController, error) {
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

	return &GetController{
		logger:    loggerService,
		builder:   videoBuilder,
		service:   videoCRUDService,
		responder: responseService,
	}, nil
}

func (c *GetController) Get(w http.ResponseWriter, r *http.Request) {
	reqDTO, err := c.builder.BuildGetRequestDTOFromRequest(r)
	if err != nil {
		c.responder.Respond(w, c.logger.LogPropagate(err))
		return
	}

	videoAgg, err := c.service.Get(reqDTO)
	if err != nil {
		c.responder.Respond(w, c.logger.LogPropagate(err))
		return
	}

	c.responder.Respond(w, videoAgg)
}

func (c *GetController) AddRoute(router *mux.Router) {
	router.
		Path(GetPath).
		HandlerFunc(c.Get).
		Methods(http.MethodGet)
}
