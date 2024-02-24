package user

import (
	builderinterface "github.com/Borislavv/video-streaming/internal/domain/builder/interface"
	"github.com/Borislavv/video-streaming/internal/domain/logger/interface"
	"github.com/Borislavv/video-streaming/internal/domain/service/di/interface"
	userinterface "github.com/Borislavv/video-streaming/internal/domain/service/user/interface"
	responseinterface "github.com/Borislavv/video-streaming/internal/infrastructure/api/v1/response/interface"
	"github.com/gorilla/mux"
	"net/http"
)

const UpdatePath = "/user/{id}"

type UpdateController struct {
	logger    loggerinterface.Logger
	builder   builderinterface.User
	service   userinterface.CRUD
	responder responseinterface.Responder
}

func NewUpdateUserController(serviceContainer diinterface.ContainerManager) (*UpdateController, error) {
	loggerService, err := serviceContainer.GetLoggerService()
	if err != nil {
		return nil, err
	}

	userBuilder, err := serviceContainer.GetUserBuilder()
	if err != nil {
		return nil, loggerService.LogPropagate(err)
	}

	userCRUDService, err := serviceContainer.GetUserCRUDService()
	if err != nil {
		return nil, loggerService.LogPropagate(err)
	}

	responseService, err := serviceContainer.GetResponderService()
	if err != nil {
		return nil, loggerService.LogPropagate(err)
	}

	return &UpdateController{
		logger:    loggerService,
		builder:   userBuilder,
		service:   userCRUDService,
		responder: responseService,
	}, nil
}

func (c *UpdateController) Update(w http.ResponseWriter, r *http.Request) {
	userReqDTO, err := c.builder.BuildUpdateRequestDTOFromRequest(r)
	if err != nil {
		c.responder.Respond(w, c.logger.LogPropagate(err))
		return
	}

	userAgg, err := c.service.Update(userReqDTO)
	if err != nil {
		c.responder.Respond(w, c.logger.LogPropagate(err))
		return
	}

	userRespDTO, err := c.builder.BuildResponseDTO(userAgg)
	if err != nil {
		c.responder.Respond(w, c.logger.LogPropagate(err))
		return
	}

	c.responder.Respond(w, userRespDTO)
}

func (c *UpdateController) AddRoute(router *mux.Router) {
	router.
		Path(UpdatePath).
		HandlerFunc(c.Update).
		Methods(http.MethodPatch)
}
