package service

import (
	"github.com/Borislavv/video-streaming/internal/domain/model"
	"github.com/Borislavv/video-streaming/internal/domain/model/stream"
	"github.com/gorilla/websocket"
)

type Manager interface {
	Get()
	Save()
	Delete()
}

type Reader interface {
	Read(resource model.Resource) chan *stream.Chunk
}

type Streamer interface {
	Stream(conn *websocket.Conn)
}
