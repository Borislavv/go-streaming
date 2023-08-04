package socket

import (
	"context"
	"github.com/Borislavv/video-streaming/internal/app/service"
	"github.com/gorilla/websocket"
	"log"
	"net"
	"net/http"
	"sync"
	"time"
)

const (
	Host = "127.0.0.1"
	Port = "9988"
	Netw = "tcp"
)

type Server struct {
	host    string
	port    string
	network string

	streamer service.Streamer

	errCh chan error
}

func NewSocketServer(streamer service.Streamer, errCh chan error) *Server {
	return &Server{
		host:     Host,
		port:     Port,
		network:  Netw,
		streamer: streamer,
		errCh:    errCh,
	}
}

// Listen is method which running a websocket server
func (s *Server) Listen(ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()

	addr, err := net.ResolveTCPAddr(s.network, net.JoinHostPort(s.host, s.port))
	if err != nil {
		s.errCh <- err
		return
	}

	server := &http.Server{
		Addr:    addr.String(),
		Handler: http.HandlerFunc(s.handleConnection),
	}

	wg.Add(1)
	go func() {
		defer wg.Done()
		defer log.Println("[socket server]: stopped")
		if lsErr := server.ListenAndServe(); lsErr != nil && lsErr != http.ErrServerClosed {
			s.errCh <- lsErr
			return
		}
	}()

	log.Println("[socket server]: running...")
	<-ctx.Done()
	log.Println("[socket server]: shutting down...")

	serverCtx, cancel := context.WithTimeout(ctx, time.Second*5)
	defer cancel()

	if sdErr := server.Shutdown(serverCtx); sdErr != nil && sdErr != context.Canceled {
		s.errCh <- sdErr
		return
	}
}

// handleConnection is method which handle each websocket connection
func (s *Server) handleConnection(w http.ResponseWriter, r *http.Request) {
	upgrader := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		s.errCh <- err
		return
	}
	defer func() {
		if err = conn.Close(); err != nil {
			s.errCh <- err
			return
		}
	}()

	log.Printf("[socket server]: accpted a new websocket connection [%s]", conn.RemoteAddr())

	s.streamer.Stream(conn)
}
