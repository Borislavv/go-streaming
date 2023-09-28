package ws

import (
	"errors"
	"fmt"
	"github.com/Borislavv/video-streaming/internal/domain/dto"
	"github.com/Borislavv/video-streaming/internal/domain/logger"
	"github.com/Borislavv/video-streaming/internal/infrastructure/service/streamer/action/enum"
	"github.com/Borislavv/video-streaming/internal/infrastructure/service/streamer/action/model"
	"github.com/gorilla/websocket"
	"strconv"
	"strings"
)

const (
	// message parts separator
	protoSeparator string = ":"
	// message prefixes
	startMsgPref string = "start"
	errMsgPref   string = "error"
	stopMsgPref  string = "stop"
)

type Communicator struct {
	logger logger.Logger
}

func NewWebSocketCommunicator(logger logger.Logger) *Communicator {
	return &Communicator{
		logger: logger,
	}
}

func (w *Communicator) Start(audioCodec string, videoCodec string, conn *websocket.Conn) error {
	b := strings.Builder{}
	b.WriteString(startMsgPref)
	b.WriteString(protoSeparator)
	b.WriteString(audioCodec)
	b.WriteString(protoSeparator)
	b.WriteString(videoCodec)
	initMessage := b.String()

	// writing the stream initialization message in a websocket connection
	if err := conn.WriteMessage(websocket.TextMessage, []byte(initMessage)); err != nil {
		return w.logger.ErrorPropagate(fmt.Sprintf("[%v]: %v", conn.RemoteAddr(), err.Error()))
	}

	return nil
}

func (w *Communicator) Send(chunk dto.Chunk, conn *websocket.Conn) error {
	if chunk.GetError() != nil {
		return w.logger.CriticalPropagate(fmt.Sprintf("[%v]: %v", conn.RemoteAddr(), chunk.GetError().Error()))
	}

	if err := conn.WriteMessage(websocket.BinaryMessage, chunk.GetData()); err != nil {
		return w.logger.CriticalPropagate(fmt.Sprintf("[%v]: %v", conn.RemoteAddr(), err.Error()))
	}

	return nil
}

func (w *Communicator) Parse(bytes []byte) (action enum.Actions, data interface{}, err error) {
	p := strings.Split(string(bytes), protoSeparator)
	if len(p) == 0 {
		return "", nil, errors.New("unable to parse message due to empty message was passed")
	}

	switch enum.Actions(p[0]) {
	case enum.StreamByID:
		return enum.Actions(p[0]), model.StreamByIdData{ID: p[1]}, nil
	case enum.StreamByIDWithOffset:
		from, ferr := strconv.ParseFloat(p[2], 64)
		if ferr != nil {
			return "", nil, ferr
		}
		duration, derr := strconv.ParseFloat(p[3], 64)
		if derr != nil {
			return "", nil, derr
		}

		return enum.Actions(p[0]), model.StreamByIdWithOffsetData{
			ID:       p[1],
			From:     from,
			Duration: duration,
		}, nil
	default:
		return "", nil, fmt.Errorf(
			"unable to parse message because received unknown action '%v'", p[0],
		)
	}
}

func (w *Communicator) Error(err error, conn *websocket.Conn) error {
	msg := []byte(fmt.Sprintf("%v:%v", errMsgPref, err.Error()))

	if e := conn.WriteMessage(websocket.TextMessage, msg); e != nil {
		return w.logger.CriticalPropagate(fmt.Sprintf("[%v]: %v", conn.RemoteAddr(), e.Error()))
	}

	return nil
}

func (w *Communicator) Stop(conn *websocket.Conn) error {
	if err := conn.WriteMessage(websocket.TextMessage, []byte(stopMsgPref)); err != nil {
		return w.logger.CriticalPropagate(fmt.Sprintf("[%v]: %v", conn.RemoteAddr(), err.Error()))
	}
	return nil
}
