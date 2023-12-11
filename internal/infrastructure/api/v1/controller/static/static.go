package static

import (
	"github.com/Borislavv/video-streaming/internal/domain/logger/interface"
	"github.com/Borislavv/video-streaming/internal/domain/service/di/interface"
	response_interface "github.com/Borislavv/video-streaming/internal/infrastructure/api/v1/response/interface"
	"github.com/Borislavv/video-streaming/internal/infrastructure/helper"
	"github.com/gorilla/mux"
	"net/http"
	"strings"
)

const ResourcesPrefix = "/static/"

type FilesController struct {
	logger    logger_interface.Logger
	responder response_interface.Responder
}

func NewFilesController(serviceContainer di_interface.ContainerManager) (*FilesController, error) {
	loggerService, err := serviceContainer.GetLoggerService()
	if err != nil {
		return nil, err
	}

	responseService, err := serviceContainer.GetResponderService()
	if err != nil {
		return nil, loggerService.LogPropagate(err)
	}

	return &FilesController{
		logger:    loggerService,
		responder: responseService,
	}, nil
}

func (c *FilesController) Serve(w http.ResponseWriter, r *http.Request) {
	dir, err := helper.StaticFilesDir()
	if err != nil {
		c.responder.Respond(w, c.logger.LogPropagate(err))
		return
	}

	path := r.URL.Path
	if strings.Contains(path, ResourcesPrefix) {
		path = strings.ReplaceAll(path, ResourcesPrefix, "")
	}

	http.ServeFile(w, r, dir+path)
}

func (c *FilesController) AddRoute(router *mux.Router) {
	router.
		PathPrefix(ResourcesPrefix).
		HandlerFunc(c.Serve).
		Methods(http.MethodGet)
}
