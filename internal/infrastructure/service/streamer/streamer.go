package streamer

import (
	"fmt"
	"github.com/Borislavv/video-streaming/internal/domain/logger"
	"github.com/Borislavv/video-streaming/internal/infrastructure/service/streamer/action/handler"
	"github.com/Borislavv/video-streaming/internal/infrastructure/service/streamer/action/listener"
	"github.com/gorilla/websocket"
	"sync"
)

type ResourceStreamer struct {
	logger   logger.Logger
	listener listener.ActionsListener
	handler  handler.ActionsHandler
}

func NewStreamingService(
	logger logger.Logger,
	listener listener.ActionsListener,
	handler handler.ActionsHandler,
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
	s.handler.Handle(wg, s.listener.Listen(wg, conn))
	wg.Wait()

	s.logger.Info(fmt.Sprintf("[%v]: streaming is stopped", conn.RemoteAddr()))
}
