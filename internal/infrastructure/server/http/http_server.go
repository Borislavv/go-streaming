package http

import (
	"context"
	"github.com/Borislavv/video-streaming/internal/infrastructure/api/v1/controller"
	"github.com/gorilla/mux"
	"log"
	"net"
	"net/http"
	"sync"
	"time"
)

const (
	Host = "127.0.0.1"
	Port = "8000"
	Netw = "tcp"

	ApiV1 = "/api/v1"
)

type Server struct {
	host        string
	port        string
	network     string
	errCh       chan error
	controllers []controller.Controller
}

func NewHttpServer(controllers []controller.Controller, errCh chan error) *Server {
	return &Server{
		host:        Host,
		port:        Port,
		network:     Netw,
		errCh:       errCh,
		controllers: controllers,
	}
}

func (s *Server) Listen(ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()
	addr, err := net.ResolveTCPAddr(s.network, net.JoinHostPort(s.host, s.port))
	if err != nil {
		s.errCh <- err
		return
	}

	server := http.Server{
		Addr:    addr.String(),
		Handler: s.addRoutes(),
	}

	wg.Add(1)
	go func() {
		defer wg.Done()
		defer log.Println("[http server]: stopped")
		if err = server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			s.errCh <- err
			return
		}
	}()

	log.Println("[http server]: running...")
	<-ctx.Done()
	log.Println("[http server]: shutting down...")

	serverCtx, cancel := context.WithTimeout(ctx, time.Second*5)
	defer cancel()

	if shErr := server.Shutdown(serverCtx); shErr != nil && shErr != context.Canceled {
		s.errCh <- shErr
		return
	}
}

func (s *Server) addRoutes() *mux.Router {
	router := mux.NewRouter()

	routerV1 := router.
		PathPrefix(ApiV1).
		Subrouter()

	for _, c := range s.controllers {
		c.AddRoute(routerV1)
	}

	return router
}
