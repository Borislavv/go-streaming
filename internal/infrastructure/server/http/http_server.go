package http

import (
	"context"
	"github.com/Borislavv/video-streaming/internal/infrastructure/api/v1/controller"
	"github.com/gorilla/mux"
	"log"
	"net"
	"net/http"
	"sync"
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
	addr, err := net.ResolveTCPAddr(s.network, net.JoinHostPort(s.host, s.port))
	if err != nil {
		s.errCh <- err
		return
	}

	if err = http.ListenAndServe(addr.String(), s.addRoutes()); err != nil {
		log.Fatalln(err)
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
