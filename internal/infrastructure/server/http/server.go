package http

import (
	"context"
	"github.com/Borislavv/video-streaming/internal/domain/api/request"
	"github.com/Borislavv/video-streaming/internal/domain/logger"
	"github.com/Borislavv/video-streaming/internal/infrastructure/api/v1/controller"
	"github.com/gorilla/mux"
	"net"
	"net/http"
	"sync"
	"time"
)

type Server struct {
	host           string // example: "0.0.0.0"
	port           string // example: "8000"
	transportProto string // example: "tcp"

	apiVersionPrefix    string // example: "/api/v1"
	renderVersionPrefix string // example: ""
	staticVersionPrefix string // example: ""

	restControllers   []controller.Controller
	renderControllers []controller.Controller
	staticControllers []controller.Controller

	logger             logger.Logger
	reqParamsExtractor request.Extractor
}

func NewHttpServer(
	host string,
	port string,
	transportProto string,
	apiVersionPrefix string,
	renderVersionPrefix string,
	staticVersionPrefix string,
	restControllers []controller.Controller,
	renderControllers []controller.Controller,
	staticControllers []controller.Controller,
	logger logger.Logger,
	reqParamsExtractor request.Extractor,
) *Server {
	return &Server{
		host:                host,
		port:                port,
		transportProto:      transportProto,
		apiVersionPrefix:    apiVersionPrefix,
		renderVersionPrefix: renderVersionPrefix,
		staticVersionPrefix: staticVersionPrefix,
		restControllers:     restControllers,
		renderControllers:   renderControllers,
		staticControllers:   staticControllers,
		logger:              logger,
		reqParamsExtractor:  reqParamsExtractor,
	}
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

	// RestAPI controllers
	restRouterV1 := router.
		PathPrefix(s.apiVersionPrefix).
		Subrouter()
	restRouterV1.
		Use(s.restApiHeaderMiddleware)
	restRouterV1.
		Use(s.requestsLoggingMiddleware)

	for _, c := range s.restControllers {
		c.AddRoute(restRouterV1)
	}

	// Native templates rendering controllers
	renderRouterV1 := router.
		PathPrefix(s.renderVersionPrefix).
		Subrouter()

	for _, c := range s.renderControllers {
		c.AddRoute(renderRouterV1)
	}

	// Static files serving controllers
	staticRouterV1 := router.
		PathPrefix(s.staticVersionPrefix).
		Subrouter()

	for _, c := range s.staticControllers {
		c.AddRoute(staticRouterV1)
	}

	return router
}

func (s *Server) restApiHeaderMiddleware(handler http.Handler) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			handler.ServeHTTP(w, r)
		},
	)
}

func (s *Server) requestsLoggingMiddleware(handler http.Handler) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			requestInfo := struct {
				Date       time.Time         `json:"date"`
				Method     string            `json:"method"`
				URL        string            `json:"URL"`
				Header     http.Header       `json:"header"`
				RemoteAddr string            `json:"remoteAddr"`
				Params     map[string]string `json:"params"`
			}{
				Method:     r.Method,
				URL:        r.URL.String(),
				Header:     r.Header,
				RemoteAddr: r.RemoteAddr,
				Params:     s.reqParamsExtractor.Parameters(r),
			}
			s.logger.LogRequestInfo(requestInfo)

			handler.ServeHTTP(w, r)
		},
	)
}
