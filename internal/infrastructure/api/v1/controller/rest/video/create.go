package video

import (
	"github.com/Borislavv/video-streaming/internal/domain/builder/interface"
	"github.com/Borislavv/video-streaming/internal/domain/logger/interface"
	authenticator_interface "github.com/Borislavv/video-streaming/internal/domain/service/authenticator/interface"
	diinterface "github.com/Borislavv/video-streaming/internal/domain/service/di/interface"
	video_interface "github.com/Borislavv/video-streaming/internal/domain/service/video/interface"
	response_interface "github.com/Borislavv/video-streaming/internal/infrastructure/api/v1/response/interface"
	"github.com/gorilla/mux"
	"net/http"
)

const CreatePath = "/video"

type CreateController struct {
	logger      loggerinterface.Logger
	builder     builder_interface.Video
	service     video_interface.CRUD
	authService authenticator_interface.Authenticator
	responder   response_interface.Responder
}

func NewCreateController(serviceContainer diinterface.ContainerManager) (*CreateController, error) {
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

	authService, err := serviceContainer.GetAuthService()
	if err != nil {
		return nil, loggerService.LogPropagate(err)
	}

	responseService, err := serviceContainer.GetResponderService()
	if err != nil {
		return nil, loggerService.LogPropagate(err)
	}

	return &CreateController{
		logger:      loggerService,
		builder:     videoBuilder,
		service:     videoCRUDService,
		authService: authService,
		responder:   responseService,
	}, nil
}

func (c *CreateController) Create(w http.ResponseWriter, r *http.Request) {
	videoDTO, err := c.builder.BuildCreateRequestDTOFromRequest(r)
	if err != nil {
		c.responder.Respond(w, c.logger.LogPropagate(err))
		return
	}

	videoAgg, err := c.service.Create(videoDTO)
	if err != nil {
		c.responder.Respond(w, c.logger.LogPropagate(err))
		return
	}

	c.responder.Respond(w, videoAgg)
	w.WriteHeader(http.StatusCreated)
}

func (c *CreateController) AddRoute(router *mux.Router) {
	router.
		Path(CreatePath).
		HandlerFunc(c.Create).
		Methods(http.MethodPost)
}
