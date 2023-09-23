package streamer

import (
	"fmt"
	"github.com/Borislavv/video-streaming/internal/domain/logger"
	"github.com/gorilla/websocket"
	"sync"
)

type WebSocketActionsListener struct {
	logger       logger.Logger
	communicator ProtoCommunicator
}

func NewWebSocketActionsListener(
	logger logger.Logger,
	proto ProtoCommunicator,
) *WebSocketActionsListener {
	return &WebSocketActionsListener{
		logger:       logger,
		communicator: proto,
	}
}

func (l *WebSocketActionsListener) Listen(wg *sync.WaitGroup, conn *websocket.Conn) <-chan Action {
	actionsCh := make(chan Action, 1)

	wg.Add(1)
	go func() {
		defer func() {
			close(actionsCh)
			wg.Done()
		}()

		for {
			t, b, err := conn.ReadMessage()
			if err != nil {
				if websocket.IsCloseError(err, websocket.CloseGoingAway, websocket.CloseNormalClosure) {
					l.logger.Info(fmt.Sprintf("[%v]: websocket connection has been closed", conn.RemoteAddr()))
					return
				}
				l.logger.Error(fmt.Sprintf("[%v]: %v", conn.RemoteAddr(), err.Error()))
				return
			}
			if t == websocket.TextMessage {
				do, data := l.communicator.Parse(b)
				if _, isSupported := supportedActionsMap[do]; isSupported {
					actionsCh <- Action{do: do, data: data}
				} else {
					l.logger.Critical(fmt.Sprintf("do: %+v, data: %+v received unsupport action", do, data))
				}
			}
		}
	}()

	return actionsCh
}
