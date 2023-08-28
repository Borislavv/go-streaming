package service

import (
	"github.com/Borislavv/video-streaming/internal/domain/entity"
	"github.com/gorilla/websocket"
)

type Manager interface {
	Get()
	Save()
	Delete()
}

type Reader interface {
	Read(resource entity.Resource) chan *entity.Chunk
}

type Streamer interface {
	Stream(conn *websocket.Conn)
}
