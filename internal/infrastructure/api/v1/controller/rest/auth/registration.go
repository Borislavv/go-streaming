package auth

import (
	"github.com/Borislavv/video-streaming/internal/domain/builder/interface"
	"github.com/Borislavv/video-streaming/internal/domain/logger/interface"
	diinterface "github.com/Borislavv/video-streaming/internal/domain/service/di/interface"
	userinterface "github.com/Borislavv/video-streaming/internal/domain/service/user/interface"
	responseinterface "github.com/Borislavv/video-streaming/internal/infrastructure/api/v1/response/interface"
	"github.com/gorilla/mux"
	"net/http"
)

const RegistrationPath = "/registration"

type RegistrationController struct {
	logger    loggerinterface.Logger
	builder   builderinterface.User
	service   userinterface.CRUD
	responder responseinterface.Responder
}

func NewRegistrationController(serviceContainer diinterface.ServiceContainer) (*RegistrationController, error) {
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

	return &RegistrationController{
		logger:    loggerService,
		builder:   userBuilder,
		service:   userCRUDService,
		responder: responseService,
	}, nil
}

// Registration - is an endpoint for create a new user.
func (c *RegistrationController) Registration(w http.ResponseWriter, r *http.Request) {
	// building a create user request DTO
	userReqDTO, err := c.builder.BuildCreateRequestDTOFromRequest(r)
	if err != nil {
		c.responder.Respond(w, c.logger.LogPropagate(err))
		return
	}

	// creating user by appropriate service
	userAgg, err := c.service.Create(userReqDTO)
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
	w.WriteHeader(http.StatusCreated)
}

func (c *RegistrationController) AddRoute(router *mux.Router) {
	router.
		Path(RegistrationPath).
		HandlerFunc(c.Registration).
		Methods(http.MethodPost)
}
