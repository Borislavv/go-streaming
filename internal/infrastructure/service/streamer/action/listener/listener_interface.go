package listener

import (
	"github.com/Borislavv/video-streaming/internal/infrastructure/service/streamer/action/model"
	"github.com/gorilla/websocket"
	"sync"
)

type ActionsListener interface {
	Listen(wg *sync.WaitGroup, conn *websocket.Conn) <-chan model.Action
}
