package auth

import (
	"github.com/Borislavv/video-streaming/internal/domain/builder"
	"github.com/Borislavv/video-streaming/internal/domain/logger"
	"github.com/Borislavv/video-streaming/internal/domain/service/user"
	"github.com/Borislavv/video-streaming/internal/infrastructure/api/v1/response"
	"github.com/gorilla/mux"
	"net/http"
)

const RegistrationPath = "/registration"

type RegistrationController struct {
	logger   logger.Logger
	builder  builder.User
	service  user.CRUD
	response response.Responder
}

func NewRegistrationController(
	logger logger.Logger,
	builder builder.User,
	service user.CRUD,
	response response.Responder,
) *RegistrationController {
	return &RegistrationController{
		logger:   logger,
		builder:  builder,
		service:  service,
		response: response,
	}
}

func (c *RegistrationController) Registration(w http.ResponseWriter, r *http.Request) {
	userDTO, err := c.builder.BuildCreateRequestDTOFromRequest(r)
	if err != nil {
		c.response.Respond(w, c.logger.LogPropagate(err))
		return
	}

	userAgg, err := c.service.Create(userDTO)
	if err != nil {
		c.response.Respond(w, c.logger.LogPropagate(err))
		return
	}

	c.response.Respond(w, userAgg)
	w.WriteHeader(http.StatusCreated)
}

func (c *RegistrationController) AddRoute(router *mux.Router) {
	router.
		Path(RegistrationPath).
		HandlerFunc(c.Registration).
		Methods(http.MethodPost)
}
