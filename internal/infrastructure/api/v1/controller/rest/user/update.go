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
	userReqDTO, err := c.builder.BuildUpdateRequestDTOFromRequest(r)
	if err != nil {
		c.response.Respond(w, c.logger.LogPropagate(err))
		return
	}

	userAgg, err := c.service.Update(userReqDTO)
	if err != nil {
		c.response.Respond(w, c.logger.LogPropagate(err))
		return
	}

	userRespDTO, err := c.builder.BuildResponseDTO(userAgg)
	if err != nil {
		c.response.Respond(w, c.logger.LogPropagate(err))
		return
	}

	c.response.Respond(w, userRespDTO)
}

func (c *UpdateController) AddRoute(router *mux.Router) {
	router.
		Path(UpdatePath).
		HandlerFunc(c.Update).
		Methods(http.MethodPatch)
}
