package http

import (
	"context"
	"github.com/Borislavv/video-streaming/internal/app/logger"
	"github.com/Borislavv/video-streaming/internal/infrastructure/api/v1/controller"
	"github.com/gorilla/mux"
	"net"
	"net/http"
	"sync"
	"time"
)

const (
	Host = "127.0.0.1"
	Port = "8000"
	Netw = "tcp"

	RestApiV1   = "/api/v1"
	RenderApiV1 = ""
	StaticApiV1 = ""
)

type Server struct {
	host              string
	port              string
	network           string
	restControllers   []controller.Controller
	renderControllers []controller.Controller
	staticControllers []controller.Controller
	logger            logger.Logger
}

func NewHttpServer(
	restControllers []controller.Controller,
	renderControllers []controller.Controller,
	staticControllers []controller.Controller,
	logger logger.Logger,
) *Server {
	return &Server{
		host:              Host,
		port:              Port,
		network:           Netw,
		restControllers:   restControllers,
		renderControllers: renderControllers,
		staticControllers: staticControllers,
		logger:            logger,
	}
}

func (s *Server) Listen(ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()
	addr, err := net.ResolveTCPAddr(s.network, net.JoinHostPort(s.host, s.port))
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
		defer s.logger.Info("[http server]: stopped")
		if err = server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			s.logger.Error(err)
			return
		}
	}()

	s.logger.Info("[http server]: running...")
	<-ctx.Done()
	s.logger.Info("[http server]: shutting down...")

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
		PathPrefix(RestApiV1).
		Subrouter()

	for _, c := range s.restControllers {
		c.AddRoute(restRouterV1)
	}

	// Native templates rendering controllers
	renderRouterV1 := router.
		PathPrefix(RenderApiV1).
		Subrouter()

	for _, c := range s.renderControllers {
		c.AddRoute(renderRouterV1)
	}

	// Static files serving controllers
	staticRouterV1 := router.
		PathPrefix(StaticApiV1).
		Subrouter()

	for _, c := range s.staticControllers {
		c.AddRoute(staticRouterV1)
	}

	return router
}
