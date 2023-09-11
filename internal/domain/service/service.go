package service

import (
	"github.com/Borislavv/video-streaming/internal/domain/agg"
	"github.com/Borislavv/video-streaming/internal/domain/dto"
	"github.com/gorilla/websocket"
)

type Resource interface {
	Upload(req dto.UploadRequest) (*agg.Resource, error)
}

type Video interface {
	Get(req dto.GetRequest) (*agg.Video, error)
	List(req dto.ListRequest) ([]*agg.Video, error)
	Create(req dto.CreateRequest) (*agg.Video, error)
	Update(req dto.UpdateRequest) (*agg.Video, error)
	Delete(req dto.DeleteRequest) error
}

type Reader interface {
	Read(resource dto.Resource) chan *dto.Chunk
}

type Streamer interface {
	Stream(conn *websocket.Conn)
}
