package service

import "github.com/gorilla/websocket"

type Streamer interface {
	Stream(conn *websocket.Conn)
}
