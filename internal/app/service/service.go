package service

import (
	"github.com/Borislavv/video-streaming/internal/domain/model"
	"github.com/gorilla/websocket"
)

type Manager interface {
	//Get()
	//Save()
	//Delete()
	Read() chan *model.Chunk
}

type Streamer interface {
	StartStreaming(connection *websocket.Conn) error
}
