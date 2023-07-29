package socket

import (
	"github.com/Borislavv/video-streaming/internal/app/service"
	"github.com/gorilla/websocket"
	"log"
	"net"
	"net/http"
)

const (
	HOST = "127.0.0.1"
	PORT = "9988"
	NETW = "tcp"
)

type Server struct {
	host    string
	port    string
	network string

	streamer service.Streamer
}

func NewSocketServer(streamer service.Streamer) *Server {
	return &Server{
		host:     HOST,
		port:     PORT,
		network:  NETW,
		streamer: streamer,
	}
}

func (s *Server) Listen() error {
	log.Println("socket server: started")
	defer log.Println("socket server: stopped")

	addr, err := net.ResolveTCPAddr(s.network, net.JoinHostPort(s.host, s.port))
	if err != nil {
		return err
	}

	http.HandleFunc("/ws", s.handleConnection)
	if err = http.ListenAndServe(addr.String(), nil); err != nil {
		return err
	}

	return nil
}

func (s *Server) handleConnection(w http.ResponseWriter, r *http.Request) {
	upgrader := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatalln(err)
	}
	defer conn.Close()

	log.Printf("socket server: accpted a new connection [%s]", conn.RemoteAddr())

	if err = s.streamer.StartStreaming(conn); err != nil {
		log.Fatalln(err)
	}
}
