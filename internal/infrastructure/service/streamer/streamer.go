package streamer

import (
	"fmt"
	"github.com/Borislavv/video-streaming/internal/domain/dto"
	"github.com/Borislavv/video-streaming/internal/domain/entity"
	"github.com/Borislavv/video-streaming/internal/domain/logger"
	"github.com/gorilla/websocket"
	"sync"
)

var (
	supportedActionsMap = map[ActionEnum]struct{}{
		streamByID: {},
	}
)

type ProtoCommunicator interface {
	Start(audioCodec string, videoCodec string, conn *websocket.Conn) error
	Send(chunk dto.Chunk, conn *websocket.Conn) error
	Parse(bytes []byte) (action ActionEnum, data string)
	Error(err error, conn *websocket.Conn) error
	Stop(conn *websocket.Conn) error
}

type CodecsDetector interface {
	Detect(resource entity.Resource) (audioCodec string, videoCodec string, err error)
}

type ActionsListener interface {
	Listen(wg *sync.WaitGroup, conn *websocket.Conn) <-chan Action
}

type ActionsHandler interface {
	Handle(wg *sync.WaitGroup, conn *websocket.Conn, actionsCh <-chan Action)
}

type ResourceStreamer struct {
	logger   logger.Logger
	listener ActionsListener
	handler  ActionsHandler
}

func NewStreamingService(
	logger logger.Logger,
	listener ActionsListener,
	handler ActionsHandler,
) *ResourceStreamer {
	return &ResourceStreamer{
		logger:   logger,
		listener: listener,
		handler:  handler,
	}
}

func (s *ResourceStreamer) HandleConn(conn *websocket.Conn) {
	s.logger.Info(fmt.Sprintf("[%v]: start streaming", conn.RemoteAddr()))

	wg := &sync.WaitGroup{}
	s.handler.Handle(wg, conn, s.listener.Listen(wg, conn))
	wg.Wait()

	s.logger.Info(fmt.Sprintf("[%v]: streaming is stopped", conn.RemoteAddr()))
}
