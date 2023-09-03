package service

import (
	"github.com/Borislavv/video-streaming/internal/domain/dto"
	"github.com/Borislavv/video-streaming/internal/domain/entity"
	"github.com/gorilla/websocket"
)

type Video interface {
	Create(video dto.CreateRequest) (string, error)
}

type Reader interface {
	Read(resource entity.Resource) chan *dto.ChunkDto
}

type Streamer interface {
	Stream(conn *websocket.Conn)
}
