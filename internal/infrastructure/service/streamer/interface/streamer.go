package streamerinterface

import "github.com/gorilla/websocket"

type Streamer interface {
	HandleConn(conn *websocket.Conn)
}
