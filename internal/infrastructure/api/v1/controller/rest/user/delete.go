package user

import (
	"github.com/Borislavv/video-streaming/internal/domain/builder"
	"github.com/Borislavv/video-streaming/internal/domain/logger"
	"github.com/Borislavv/video-streaming/internal/domain/service/user"
	"github.com/Borislavv/video-streaming/internal/infrastructure/api/v1/response"
	"github.com/gorilla/mux"
	"net/http"
)

const DeletePath = "/user/{id}"

type DeleteUserController struct {
	logger   logger.Logger
	builder  builder.User
	service  user.CRUD
	response response.Responder
}

func NewDeleteController(
	logger logger.Logger,
	builder builder.User,
	service user.CRUD,
	response response.Responder,
) *DeleteUserController {
	return &DeleteUserController{
		logger:   logger,
		builder:  builder,
		service:  service,
		response: response,
	}
}

func (c *DeleteUserController) Delete(w http.ResponseWriter, r *http.Request) {
	reqDTO, err := c.builder.BuildDeleteRequestDTOFromRequest(r)
	if err != nil {
		c.response.Respond(w, c.logger.LogPropagate(err))
		return
	}

	if err = c.service.Delete(reqDTO); err != nil {
		c.response.Respond(w, c.logger.LogPropagate(err))
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (c *DeleteUserController) AddRoute(router *mux.Router) {
	router.
		Path(DeletePath).
		HandlerFunc(c.Delete).
		Methods(http.MethodDelete)
}
