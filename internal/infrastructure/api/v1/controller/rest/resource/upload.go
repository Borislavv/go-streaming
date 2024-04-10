package resource

import (
	builderinterface "github.com/Borislavv/video-streaming/internal/domain/builder/interface"
	"github.com/Borislavv/video-streaming/internal/domain/logger/interface"
	"github.com/Borislavv/video-streaming/internal/domain/service/di/interface"
	resourceinterface "github.com/Borislavv/video-streaming/internal/domain/service/resource/interface"
	responseinterface "github.com/Borislavv/video-streaming/internal/infrastructure/api/v1/response/interface"
	"github.com/gorilla/mux"
	"net/http"
)

const UploadPath = "/resource"

type UploadResourceController struct {
	logger    loggerinterface.Logger
	builder   builderinterface.Resource
	service   resourceinterface.CRUD
	responder responseinterface.Responder
}

func NewUploadController(serviceContainer diinterface.ServiceContainer) (*UploadResourceController, error) {
	loggerService, err := serviceContainer.GetLoggerService()
	if err != nil {
		return nil, err
	}

	resourceBuilder, err := serviceContainer.GetResourceBuilder()
	if err != nil {
		return nil, loggerService.LogPropagate(err)
	}

	resourceCRUDService, err := serviceContainer.GetResourceCRUDService()
	if err != nil {
		return nil, loggerService.LogPropagate(err)
	}

	responseService, err := serviceContainer.GetResponderService()
	if err != nil {
		return nil, loggerService.LogPropagate(err)
	}

	return &UploadResourceController{
		logger:    loggerService,
		builder:   resourceBuilder,
		service:   resourceCRUDService,
		responder: responseService,
	}, nil
}

func (c *UploadResourceController) Upload(w http.ResponseWriter, r *http.Request) {
	reqDTO, err := c.builder.BuildUploadRequestDTOFromRequest(r)
	if err != nil {
		c.responder.Respond(w, c.logger.LogPropagate(err))
		return
	}

	resourceAgg, err := c.service.Upload(reqDTO)
	if err != nil {
		c.responder.Respond(w, c.logger.LogPropagate(err))
		return
	}

	c.responder.Respond(w, resourceAgg)
}

func (c *UploadResourceController) AddRoute(router *mux.Router) {
	router.
		Path(UploadPath).
		HandlerFunc(c.Upload).
		Methods(http.MethodPost)
}
