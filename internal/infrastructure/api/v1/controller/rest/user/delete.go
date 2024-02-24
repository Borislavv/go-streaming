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

const DeletePath = "/user/{id}"

type DeleteController struct {
	logger    loggerinterface.Logger
	builder   builderinterface.User
	service   userinterface.CRUD
	responder responseinterface.Responder
}

func NewDeleteController(serviceContainer diinterface.ContainerManager) (*DeleteController, error) {
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

	return &DeleteController{
		logger:    loggerService,
		builder:   userBuilder,
		service:   userCRUDService,
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
