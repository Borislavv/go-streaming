package http

const (
	HOST = "127.0.0.1"
	PORT = "8080"
	NETW = "tcp"
)

type Server struct {
	host    string
	port    string
	network string
}

func NewHttpServer() *Server {
	return &Server{
		host:    HOST,
		port:    PORT,
		network: NETW,
	}
}

func (s *Server) Listen() error {

}
