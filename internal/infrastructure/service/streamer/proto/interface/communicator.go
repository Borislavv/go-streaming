package proto

import (
	"github.com/Borislavv/video-streaming/internal/infrastructure/service/streamer/action/enum"
	"github.com/gorilla/websocket"
)

type Communicator interface {
	Start(audioCodec string, videoCodec string, conn *websocket.Conn) error
	Send(chunk dto_interface.Chunk, conn *websocket.Conn) error
	Parse(bytes []byte) (action enum.Actions, data interface{}, err error)
	Error(err error, conn *websocket.Conn) error
	Stop(conn *websocket.Conn) error
}
