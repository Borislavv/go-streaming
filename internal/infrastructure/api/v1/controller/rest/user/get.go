package user

import (
	"encoding/json"
	"github.com/Borislavv/video-streaming/internal/domain/agg"
	"github.com/Borislavv/video-streaming/internal/domain/builder/interface"
	"github.com/Borislavv/video-streaming/internal/domain/dto"
	"github.com/Borislavv/video-streaming/internal/domain/errors"
	"github.com/Borislavv/video-streaming/internal/domain/logger/interface"
	cacher_interface "github.com/Borislavv/video-streaming/internal/domain/service/cacher/interface"
	"github.com/Borislavv/video-streaming/internal/domain/service/di/interface"
	user_interface "github.com/Borislavv/video-streaming/internal/domain/service/user/interface"
	response_interface "github.com/Borislavv/video-streaming/internal/infrastructure/api/v1/response/interface"
	"github.com/Borislavv/video-streaming/internal/infrastructure/helper"
	"github.com/gorilla/mux"
	"net/http"
	"reflect"
	"time"
)

const (
	GetPath  = "/user"
	cacheTTL = time.Minute * 60
)

type GetController struct {
	logger    logger_interface.Logger
	builder   builder_interface.User
	service   user_interface.CRUD
	cacher    cacher_interface.Cacher
	responder response_interface.Responder
}

func NewGetController(serviceContainer di_interface.ContainerManager) (*GetController, error) {
	loggerService, err := serviceContainer.GetLoggerService()
	if err != nil {
		return nil, err
	}

	cacheService, err := serviceContainer.GetCacheService()
	if err != nil {
		return nil, loggerService.LogPropagate(err)
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

	return &GetController{
		logger:    loggerService,
		builder:   userBuilder,
		service:   userCRUDService,
		responder: responseService,
		cacher:    cacheService,
	}, nil
}

func (c *GetController) Get(w http.ResponseWriter, r *http.Request) {
	userReqDTO, err := c.builder.BuildGetRequestDTOFromRequest(r)
	if err != nil {
		c.responder.Respond(w, c.logger.LogPropagate(err))
		return
	}

	userAgg, err := c.getCached(userReqDTO)
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

func (c *GetController) getCached(reqDTO *dto.UserGetRequestDTO) (*agg.User, error) {
	key, err := json.Marshal(reqDTO)
	if err != nil {
		return nil, c.logger.LogPropagate(err)
	}

	cacheKey := helper.MD5(key)

	data, err := c.cacher.Get(
		cacheKey,
		func(item cacher_interface.CacheItem) (data interface{}, err error) {
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

func (c *GetController) AddRoute(router *mux.Router) {
	router.
		Path(GetPath).
		HandlerFunc(c.Get).
		Methods(http.MethodGet)
}
