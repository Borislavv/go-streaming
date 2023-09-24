package model

import (
	"github.com/Borislavv/video-streaming/internal/infrastructure/service/streamer/action/enum"
	"github.com/gorilla/websocket"
)

type Action struct {
	Do   enum.Actions
	Data string
	Conn *websocket.Conn
}
