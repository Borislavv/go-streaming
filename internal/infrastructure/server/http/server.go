package http

import (
	"context"
	"net"
	"net/http"
	"sync"
	"time"

	"github.com/Borislavv/video-streaming/internal/domain/enum"
	logger_interface "github.com/Borislavv/video-streaming/internal/domain/logger/interface"
	authenticator_interface "github.com/Borislavv/video-streaming/internal/domain/service/authenticator/interface"
	di_interface "github.com/Borislavv/video-streaming/internal/domain/service/di/interface"
	extractor_interface "github.com/Borislavv/video-streaming/internal/domain/service/extractor/interface"
	"github.com/Borislavv/video-streaming/internal/infrastructure/api/v1/controller"
	"github.com/Borislavv/video-streaming/internal/infrastructure/api/v1/controller/render"
	"github.com/Borislavv/video-streaming/internal/infrastructure/api/v1/request"
	response_interface "github.com/Borislavv/video-streaming/internal/infrastructure/api/v1/response/interface"
	"github.com/Borislavv/video-streaming/internal/infrastructure/helper/ruid"
	"github.com/gorilla/mux"
)

type Server struct {
	ctx context.Context

	host           string // example: "0.0.0.0"
	port           string // example: "8000"
	transportProto string // example: "tcp"

	apiVersionPrefix    string // example: "/api/v1"
	renderVersionPrefix string // example: ""
	staticVersionPrefix string // example: ""

	restAuthedControllers     []controller.Controller
	restUnauthedControllers   []controller.Controller
	renderAuthedControllers   []controller.Controller
	renderUnauthedControllers []controller.Controller
	staticControllers         []controller.Controller

	logger             logger_interface.Logger
	authService        authenticator_interface.Authenticator
	reqParamsExtractor extractor_interface.RequestParams
	responder          response_interface.Responder
}

func NewHttpServer(
	serviceContainer di_interface.ContainerManager,
	restAuthedControllers []controller.Controller,
	restUnauthedControllers []controller.Controller,
	renderAuthedControllers []controller.Controller,
	renderUnauthedControllers []controller.Controller,
	staticControllers []controller.Controller,
) (*Server, error) {
	loggerService, err := serviceContainer.GetLoggerService()
	if err != nil {
		return nil, err
	}

	ctx, err := serviceContainer.GetCtx()
	if err != nil {
		return nil, loggerService.LogPropagate(err)
	}

	authService, err := serviceContainer.GetAuthService()
	if err != nil {
		return nil, loggerService.LogPropagate(err)
	}

	requestParametersExtractorService, err := serviceContainer.GetRequestParametersExtractorService()
	if err != nil {
		return nil, loggerService.LogPropagate(err)
	}

	responderService, err := serviceContainer.GetResponderService()
	if err != nil {
		return nil, loggerService.LogPropagate(err)
	}

	cfg, err := serviceContainer.GetConfig()
	if err != nil {
		return nil, loggerService.LogPropagate(err)
	}

	return &Server{
		ctx:                       ctx,
		host:                      cfg.ResourcesHost,
		port:                      cfg.ResourcesPort,
		transportProto:            cfg.ResourcesTransport,
		apiVersionPrefix:          cfg.ResourcesApiVersionPrefix,
		renderVersionPrefix:       cfg.ResourcesRenderVersionPrefix,
		staticVersionPrefix:       cfg.ResourcesStaticVersionPrefix,
		restAuthedControllers:     restAuthedControllers,
		restUnauthedControllers:   restUnauthedControllers,
		renderAuthedControllers:   renderAuthedControllers,
		renderUnauthedControllers: renderUnauthedControllers,
		staticControllers:         staticControllers,
		logger:                    loggerService,
		authService:               authService,
		reqParamsExtractor:        requestParametersExtractorService,
		responder:                 responderService,
	}, nil
}

func (s *Server) Listen(ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()
	addr, err := net.ResolveTCPAddr(s.transportProto, net.JoinHostPort(s.host, s.port))
	if err != nil {
		s.logger.Error(err)
		return
	}

	server := http.Server{
		Addr:    addr.String(),
		Handler: s.addRoutes(),
	}

	wg.Add(1)
	go func() {
		defer wg.Done()
		defer s.logger.Info("stopped")
		if err = server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			s.logger.Error(err)
			return
		}
	}()

	s.logger.Info("running...")
	<-ctx.Done()
	s.logger.Info("shutting down...")

	serverCtx, cancel := context.WithTimeout(ctx, time.Second*5)
	defer cancel()

	if shErr := server.Shutdown(serverCtx); shErr != nil && shErr != context.Canceled {
		s.logger.Critical(err)
		return
	}
}

func (s *Server) addRoutes() *mux.Router {
	router := mux.NewRouter()

	// [AUTHED] rest api controllers which requires authorization token
	restAuthedRouterV1 := router.
		PathPrefix(s.apiVersionPrefix).
		Subrouter()
	restAuthedRouterV1.
		Use(
			s.restApiHeaderMiddleware,
			s.requestsLoggingMiddleware,
			s.restAuthorizationMiddleware,
		)

	for _, c := range s.restAuthedControllers {
		c.AddRoute(restAuthedRouterV1)
	}

	// [UNAUTHED] rest api controllers which is not requires authorization token
	restUnauthedRouterV1 := router.
		PathPrefix(s.apiVersionPrefix).
		Subrouter()
	restUnauthedRouterV1.
		Use(
			s.requestsLoggingMiddleware,
			s.restApiHeaderMiddleware,
		)

	for _, c := range s.restUnauthedControllers {
		c.AddRoute(restUnauthedRouterV1)
	}

	// [AUTHED] native templates rendering controllers
	renderAuthedRouterV1 := router.
		PathPrefix(s.renderVersionPrefix).
		Subrouter()
	renderAuthedRouterV1.
		Use(
			s.requestsLoggingMiddleware,
			s.renderAuthorizationMiddleware,
		)

	for _, c := range s.renderAuthedControllers {
		c.AddRoute(renderAuthedRouterV1)
	}

	// [UNAUTHED] native templates rendering controllers
	renderUnauthedRouterV1 := router.
		PathPrefix(s.renderVersionPrefix).
		Subrouter()
	renderUnauthedRouterV1.
		Use(
			s.requestsLoggingMiddleware,
		)

	for _, c := range s.renderUnauthedControllers {
		c.AddRoute(renderUnauthedRouterV1)
	}

	// static files serving controllers
	staticRouterV1 := router.
		PathPrefix(s.staticVersionPrefix).
		Subrouter()

	for _, c := range s.staticControllers {
		c.AddRoute(staticRouterV1)
	}

	return router
}

// restAuthorizationMiddleware checks whether the access token was passed
// and user is authed otherwise throws access denied error.
func (s *Server) restAuthorizationMiddleware(handler http.Handler) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			userID, err := s.authService.IsAuthed(r)
			if err != nil {
				s.responder.Respond(w, s.logger.LogPropagate(err))
				return
			}
			// create a new context with userID value
			ctx := context.WithValue(r.Context(), enum.UserIDContextKey, userID)
			// serve the next layer
			handler.ServeHTTP(w, r.WithContext(ctx))
		},
	)
}

// renderAuthorizationMiddleware checks whether the access token was passed
// and user is authed otherwise redirects to login page.
func (s *Server) renderAuthorizationMiddleware(handler http.Handler) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			userID, err := s.authService.IsAuthed(r)
			if err != nil {
				// error logging
				s.logger.Log(err)
				// info action logging
				s.logger.Info("redirect to login page")
				// redirecting a client to the login page
				http.Redirect(w, r, render.LoginPath, http.StatusSeeOther)
				return
			}
			// create a new context with userID value
			ctx := context.WithValue(r.Context(), enum.UserIDContextKey, userID)
			// serve the next layer
			handler.ServeHTTP(w, r.WithContext(ctx))
		},
	)
}

func (s *Server) restApiHeaderMiddleware(handler http.Handler) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			// adding the rest api header
			w.Header().Set("Content-Type", "application/json")
			// serve the next layer
			handler.ServeHTTP(w, r)
		},
	)
}

func (s *Server) requestsLoggingMiddleware(handler http.Handler) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			uniqueReqID := ruid.RequestUniqueID(r)
			requestData := &request.LoggableData{
				Date:       time.Now(),
				ReqID:      uniqueReqID,
				Type:       request.LogType,
				Method:     r.Method,
				URL:        r.URL.String(),
				Header:     r.Header,
				RemoteAddr: r.RemoteAddr,
				Params:     s.reqParamsExtractor.Parameters(r),
			}
			// request logging
			s.logger.LogData(requestData)
			// pass a requestID through entire app.
			s.logger.SetContext(context.WithValue(s.ctx, enum.UniqueRequestIDKey, uniqueReqID))
			// serve the next layer
			handler.ServeHTTP(w, r)
		},
	)
}
