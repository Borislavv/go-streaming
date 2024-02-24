package video

import (
	"github.com/Borislavv/video-streaming/internal/domain/builder/interface"
	"github.com/Borislavv/video-streaming/internal/domain/logger/interface"
	"github.com/Borislavv/video-streaming/internal/domain/service/di/interface"
	videointerface "github.com/Borislavv/video-streaming/internal/domain/service/video/interface"
	responseinterface "github.com/Borislavv/video-streaming/internal/infrastructure/api/v1/response/interface"
	"github.com/gorilla/mux"
	"net/http"
)

const ListPath = "/video"

type ListController struct {
	logger    loggerinterface.Logger
	builder   builderinterface.Video
	service   videointerface.CRUD
	responder responseinterface.Responder
}

func NewListController(serviceContainer diinterface.ContainerManager) (*ListController, error) {
	loggerService, err := serviceContainer.GetLoggerService()
	if err != nil {
		return nil, err
	}

	videoBuilder, err := serviceContainer.GetVideoBuilder()
	if err != nil {
		return nil, loggerService.LogPropagate(err)
	}

	videoCRUDService, err := serviceContainer.GetVideoCRUDService()
	if err != nil {
		return nil, loggerService.LogPropagate(err)
	}

	responseService, err := serviceContainer.GetResponderService()
	if err != nil {
		return nil, loggerService.LogPropagate(err)
	}

	return &ListController{
		logger:    loggerService,
		builder:   videoBuilder,
		service:   videoCRUDService,
		responder: responseService,
	}, nil
}

func (c *ListController) List(w http.ResponseWriter, r *http.Request) {
	reqDTO, e := c.builder.BuildListRequestDTOFromRequest(r)
	if e != nil {
		c.responder.Respond(w, c.logger.LogPropagate(e))
		return
	}

	aggList, total, err := c.service.List(reqDTO)
	if err != nil {
		c.responder.Respond(w, c.logger.LogPropagate(err))
		return
	}

	// TODO must be refactored to paginated list DTO.
	c.responder.Respond(w,
		map[string]interface{}{
			"list": aggList,
			"pagination": map[string]interface{}{
				"page":  reqDTO.Page,
				"limit": reqDTO.Limit,
				"total": total,
			},
		},
	)
}

func (c *ListController) AddRoute(router *mux.Router) {
	router.
		Path(ListPath).
		HandlerFunc(c.List).
		Methods(http.MethodGet)
}
