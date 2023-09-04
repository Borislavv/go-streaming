package service

import (
	"github.com/Borislavv/video-streaming/internal/domain/dto"
	"github.com/Borislavv/video-streaming/internal/domain/entity"
	"github.com/Borislavv/video-streaming/internal/domain/vo"
	"github.com/gorilla/websocket"
)

type Video interface {
	Create(video dto.CreateRequest) (*vo.ID, error)
}

type Reader interface {
	Read(resource entity.Resource) chan *dto.ChunkDto
}

type Streamer interface {
	Stream(conn *websocket.Conn)
}
