package user

import (
	"github.com/Borislavv/video-streaming/internal/domain/builder"
	"github.com/Borislavv/video-streaming/internal/domain/logger"
	"github.com/Borislavv/video-streaming/internal/domain/service/user"
	"github.com/Borislavv/video-streaming/internal/infrastructure/api/v1/response"
	"github.com/gorilla/mux"
	"net/http"
)

const GetPath = "/user/{id}"

// GetUserController - not implemented yet.
type GetUserController struct {
	logger   logger.Logger
	builder  builder.User
	service  user.CRUD
	response response.Responder
}

func NewGetController(
	logger logger.Logger,
	builder builder.User,
	service user.CRUD,
	response response.Responder,
) *GetUserController {
	return &GetUserController{
		logger:   logger,
		builder:  builder,
		service:  service,
		response: response,
	}
}

func (c *GetUserController) Get(w http.ResponseWriter, r *http.Request) {
	reqDTO, err := c.builder.BuildGetRequestDTOFromRequest(r)
	if err != nil {
		c.response.Respond(w, c.logger.LogPropagate(err))
		return
	}

	userAgg, err := c.service.Get(reqDTO)
	if err != nil {
		c.response.Respond(w, c.logger.LogPropagate(err))
		return
	}

	c.response.Respond(w, userAgg)
}

func (c *GetUserController) AddRoute(router *mux.Router) {
	router.
		Path(GetPath).
		HandlerFunc(c.Get).
		Methods(http.MethodGet)
}
