package streamer

import (
	"context"
	"fmt"
	"github.com/Borislavv/video-streaming/internal/domain/dto"
	"github.com/Borislavv/video-streaming/internal/domain/entity"
	"github.com/Borislavv/video-streaming/internal/domain/errs"
	"github.com/Borislavv/video-streaming/internal/domain/logger"
	"github.com/Borislavv/video-streaming/internal/domain/repository"
	"github.com/Borislavv/video-streaming/internal/domain/service"
	"github.com/Borislavv/video-streaming/internal/domain/vo"
	"github.com/gorilla/websocket"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"log"
	"sync"
)

const (
	StreamByID        ActionEnum = "ID"
	DecreaseBufferCap ActionEnum = "decrBuff"
)

type Action struct {
	do   ActionEnum
	data string
}

type ActionEnum string

func (a ActionEnum) String() string {
	return string(a)
}

var (
	supportedActionsMap = map[ActionEnum]struct{}{
		StreamByID:        {},
		DecreaseBufferCap: {},
	}
)

type ProtoCommunicator interface {
	Start(audioCodec string, videoCodec string, conn *websocket.Conn) error
	Send(chunk dto.Chunk, conn *websocket.Conn) error
	Parse(bytes []byte) (action ActionEnum, data string)
	Error(err error, conn *websocket.Conn) error
	Stop(conn *websocket.Conn) error
}

type CodecsDeterminer interface {
	Determine(resource entity.Resource) (audioCodec string, videoCodec string, err error)
}

type ResourceStreamer struct {
	ctx             context.Context
	logger          logger.Logger
	reader          service.Reader
	videoRepository repository.Video
	proto           ProtoCommunicator
	codecs          CodecsDeterminer
}

func NewStreamingService(
	ctx context.Context,
	logger logger.Logger,
	reader service.Reader,
	videoRepository repository.Video,
	wsProto ProtoCommunicator,
	resourceCodecs CodecsDeterminer,
) *ResourceStreamer {
	return &ResourceStreamer{
		ctx:             ctx,
		reader:          reader,
		logger:          logger,
		videoRepository: videoRepository,
		proto:           wsProto,
		codecs:          resourceCodecs,
	}
}

func (s *ResourceStreamer) Stream(conn *websocket.Conn) {
	s.logger.Info(fmt.Sprintf("[%v]: start streaming", conn.RemoteAddr()))

	actionCh := make(chan Action, 1)
	decrBuffCh := make(chan struct{})

	wg := &sync.WaitGroup{}
	wg.Add(2)
	go s.listenClient(wg, conn, actionCh)
	go s.handleStreamActions(wg, conn, actionCh, decrBuffCh)
	wg.Wait()

	s.logger.Info(fmt.Sprintf("[%v]: streaming is stopped", conn.RemoteAddr()))
}

func (s *ResourceStreamer) listenClient(wg *sync.WaitGroup, conn *websocket.Conn, actionsCh chan<- Action) {
	defer func() {
		close(actionsCh)
		wg.Done()
	}()

	for {
		t, b, err := conn.ReadMessage()
		if err != nil {
			if websocket.IsCloseError(err, websocket.CloseGoingAway, websocket.CloseNormalClosure) {
				s.logger.Info(fmt.Sprintf("[%v]: websocket connection has been closed", conn.RemoteAddr()))
				return
			}
			s.logger.Error(fmt.Sprintf("[%v]: %v", conn.RemoteAddr(), err.Error()))
			return
		}
		if t == websocket.TextMessage {
			do, data := s.proto.Parse(b)
			if _, isSupported := supportedActionsMap[do]; isSupported {
				actionsCh <- Action{do: do, data: data}
			} else {
				s.logger.Critical(fmt.Sprintf("do: %+v, data: %+v received unsupport action", do, data))
			}
		}
	}
}

func (s *ResourceStreamer) handleStreamActions(
	wg *sync.WaitGroup,
	conn *websocket.Conn,
	actionCh <-chan Action,
	decrBuffCapCh chan<- struct{},
) {
	defer func() {
		close(decrBuffCapCh)
		wg.Done()
	}()

	for action := range actionCh {
		// todo must be moved in the strategies of actions
		switch action.do {
		case StreamByID:
			oid, err := primitive.ObjectIDFromHex(action.data)
			if err != nil {
				s.logger.Log(err)
				continue
			}
			v, err := s.videoRepository.Find(s.ctx, vo.ID{Value: oid})
			if err != nil {
				s.logger.Log(err)
				if errs.IsNotFoundError(err) {
					if err = s.proto.Error(err, conn); err != nil {
						s.logger.Log(err)
						log.Println(" -------------------------------====== NOT FOUND ERROR (ERROR) -------------------------------======")
					}
					log.Println(" -------------------------------====== NOT FOUND ERROR -------------------------------======")
					continue
				}
				continue
			}
			s.logger.Info(fmt.Sprintf("[%v]: streaming 'resource':'%v'", conn.RemoteAddr(), v.Resource.Name))
			s.streamResource(v.Resource, conn)
		case DecreaseBufferCap:
			decrBuffCapCh <- struct{}{}
			continue
		}
	}
}

func (s *ResourceStreamer) streamResource(resource entity.Resource, conn *websocket.Conn) {
	audioCodec, videoCodec, err := s.codecs.Determine(resource)
	if err != nil {
		s.logger.Error(fmt.Sprintf("[%v]: %v", conn.RemoteAddr(), err.Error()))
		return
	}

	if err = s.proto.Start(audioCodec, videoCodec, conn); err != nil {
		s.logger.Critical(fmt.Sprintf("[%v]: %v", conn.RemoteAddr(), err.Error()))
		return
	}

	for chunk := range s.reader.Read(resource) {
		if err = s.proto.Send(chunk, conn); err != nil {
			s.logger.Critical(fmt.Sprintf("[%v]: %v", conn.RemoteAddr(), err.Error()))
			break
		}

		s.logger.Info(
			fmt.Sprintf(
				"[%v]: wrote %d bytes of '%v' to websocket",
				conn.RemoteAddr(), chunk.GetLen(), resource.Name,
			),
		)
	}

	if err = s.proto.Stop(conn); err != nil {
		s.logger.Critical(fmt.Sprintf("[%v]: %v", conn.RemoteAddr(), err.Error()))
		return
	}
}
