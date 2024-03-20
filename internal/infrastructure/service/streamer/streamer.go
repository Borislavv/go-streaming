package streamer

import (
	"fmt"
	"github.com/Borislavv/video-streaming/internal/domain/logger/interface"
	"github.com/Borislavv/video-streaming/internal/domain/service/di/interface"
	handlerinterface "github.com/Borislavv/video-streaming/internal/infrastructure/service/streamer/action/handler/interface"
	listenerinterface "github.com/Borislavv/video-streaming/internal/infrastructure/service/streamer/action/listener/interface"
	"github.com/gorilla/websocket"
	"sync"
)

type ResourceStreamer struct {
	logger   loggerinterface.Logger
	listener listenerinterface.ActionsListener
	handler  handlerinterface.ActionsHandler
}

func NewStreamingService(serviceContainer diinterface.ContainerManager) (*ResourceStreamer, error) {
	loggerService, err := serviceContainer.GetLoggerService()
	if err != nil {
		return nil, err
	}

	webSocketListener, err := serviceContainer.GetWebSocketListener()
	if err != nil {
		return nil, loggerService.LogPropagate(err)
	}

	webSocketHandler, err := serviceContainer.GetWebSocketHandler()
	if err != nil {
		return nil, loggerService.LogPropagate(err)
	}

	return &ResourceStreamer{
		logger:   loggerService,
		listener: webSocketListener,
		handler:  webSocketHandler,
	}, nil
}

func (s *ResourceStreamer) HandleConn(conn *websocket.Conn) {
	s.logger.Info(fmt.Sprintf("[%v]: start streaming", conn.RemoteAddr()))

	wg := &sync.WaitGroup{}
	s.handler.Handle(wg, s.listener.Listen(wg, conn))
	wg.Wait()

	s.logger.Info(fmt.Sprintf("[%v]: streaming is stopped", conn.RemoteAddr()))
}
