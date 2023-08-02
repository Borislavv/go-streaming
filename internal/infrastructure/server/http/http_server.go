package http

import (
	"context"
	"github.com/Borislavv/video-streaming/internal/infrastructure/api/v1/controller/video"
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
	host    string
	port    string
	network string
	errCh   chan error
}

func NewHttpServer(errCh chan error) *Server {
	return &Server{
		host:    Host,
		port:    Port,
		network: Netw,
		errCh:   errCh,
	}
}

func (s *Server) Listen(ctx context.Context, wg *sync.WaitGroup) {
	addr, err := net.ResolveTCPAddr(s.network, net.JoinHostPort(s.host, s.port))
	if err != nil {
		s.errCh <- err
		return
	}

	router := s.addRoutes()

	router.PathPrefix(ApiV1)

	//http.Handle(ApiV1, router)

	if err = http.ListenAndServe(addr.String(), router); err != nil {
		log.Fatalln(err)
	}
}

func (s *Server) addRoutes() *mux.Router {
	router := mux.NewRouter()

	// Video
	router.
		HandleFunc(video.CreatePath, video.Create).
		Methods(http.MethodPost)

	router.
		HandleFunc(video.GetPath, video.Get).
		Methods(http.MethodGet)

	router.
		HandleFunc(video.ListPath, video.List).
		Methods(http.MethodGet)

	router.
		HandleFunc(video.UpdatePath, video.Update).
		Methods(http.MethodPatch)

	// Audio
	// todo not implemented yet

	return router
}
