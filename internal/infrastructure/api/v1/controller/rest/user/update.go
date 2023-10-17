package user

import (
	"github.com/Borislavv/video-streaming/internal/domain/builder"
	"github.com/Borislavv/video-streaming/internal/domain/logger"
	"github.com/Borislavv/video-streaming/internal/domain/service/user"
	"github.com/Borislavv/video-streaming/internal/infrastructure/api/v1/response"
	"github.com/gorilla/mux"
	"net/http"
)

const UpdatePath = "/user/{id}"

type UpdateController struct {
	logger   logger.Logger
	builder  builder.User
	service  user.CRUD
	response response.Responder
}

func NewUpdateUserController(
	logger logger.Logger,
	builder builder.User,
	service user.CRUD,
	response response.Responder,
) *UpdateController {
	return &UpdateController{
		logger:   logger,
		builder:  builder,
		service:  service,
		response: response,
	}
}

func (c *UpdateController) Update(w http.ResponseWriter, r *http.Request) {
	userDTO, err := c.builder.BuildUpdateRequestDTOFromRequest(r)
	if err != nil {
		c.response.Respond(w, c.logger.LogPropagate(err))
		return
	}

	userAgg, err := c.service.Update(userDTO)
	if err != nil {
		c.response.Respond(w, c.logger.LogPropagate(err))
		return
	}

	c.response.Respond(w, userAgg)
}

func (c *UpdateController) AddRoute(router *mux.Router) {
	router.
		Path(UpdatePath).
		HandlerFunc(c.Update).
		Methods(http.MethodPatch)
}
