package user

import (
	"encoding/json"
	"github.com/Borislavv/video-streaming/internal/domain/agg"
	"github.com/Borislavv/video-streaming/internal/domain/builder"
	"github.com/Borislavv/video-streaming/internal/domain/dto"
	"github.com/Borislavv/video-streaming/internal/domain/errors"
	"github.com/Borislavv/video-streaming/internal/domain/logger"
	"github.com/Borislavv/video-streaming/internal/domain/service/user"
	"github.com/Borislavv/video-streaming/internal/infrastructure/api/v1/response"
	"github.com/Borislavv/video-streaming/internal/infrastructure/helper"
	"github.com/Borislavv/video-streaming/internal/infrastructure/service/cacher"
	"github.com/gorilla/mux"
	"net/http"
	"reflect"
	"time"
)

const (
	GetPath  = "/user/{id}"
	cacheTTL = time.Minute * 60
)

type GetUserController struct {
	logger   logger.Logger
	builder  builder.User
	service  user.CRUD
	cacher   cacher.Cacher
	response response.Responder
}

func NewGetController(
	logger logger.Logger,
	builder builder.User,
	service user.CRUD,
	cacher cacher.Cacher,
	response response.Responder,
) *GetUserController {
	return &GetUserController{
		logger:   logger,
		builder:  builder,
		service:  service,
		cacher:   cacher,
		response: response,
	}
}

func (c *GetUserController) Get(w http.ResponseWriter, r *http.Request) {
	reqDTO, err := c.builder.BuildGetRequestDTOFromRequest(r)
	if err != nil {
		c.response.Respond(w, c.logger.LogPropagate(err))
		return
	}

	userAgg, err := c.getCached(reqDTO)
	if err != nil {
		c.response.Respond(w, c.logger.LogPropagate(err))
		return
	}

	c.response.Respond(w, userAgg)
}

func (c *GetUserController) getCached(reqDTO *dto.UserGetRequestDTO) (*agg.User, error) {
	key, err := json.Marshal(reqDTO)
	if err != nil {
		return nil, c.logger.LogPropagate(err)
	}

	cacheKey := helper.MD5(key)

	data, err := c.cacher.Get(
		cacheKey,
		func(item cacher.CacheItem) (data interface{}, err error) {
			item.SetTTL(cacheTTL)
			return c.service.Get(reqDTO)
		},
	)
	if err != nil {
		return nil, c.logger.LogPropagate(err)
	}

	userAgg, ok := data.(*agg.User)
	if !ok {
		return nil, errors.NewCachedDataTypeWasNotMatchedError(
			cacheKey, reflect.TypeOf(&agg.User{}), reflect.TypeOf(data),
		)
	}

	return userAgg, err
}

func (c *GetUserController) AddRoute(router *mux.Router) {
	router.
		Path(GetPath).
		HandlerFunc(c.Get).
		Methods(http.MethodGet)
}
