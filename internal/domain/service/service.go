package service

import (
	"github.com/Borislavv/video-streaming/internal/domain/agg"
	"github.com/Borislavv/video-streaming/internal/domain/dto"
	"github.com/Borislavv/video-streaming/internal/domain/entity"
	"github.com/gorilla/websocket"
)

type Video interface {
	Get(req dto.GetRequest) (*agg.Video, error)
	List(req dto.ListRequest) ([]*agg.Video, error)
	Create(req dto.CreateRequest) (*agg.Video, error)
	Update(req dto.UpdateRequest) (*agg.Video, error)
	Delete(req dto.DeleteRequest) error
}

type Reader interface {
	Read(resource entity.Resource) chan *dto.ChunkDto
}

type Streamer interface {
	Stream(conn *websocket.Conn)
}
