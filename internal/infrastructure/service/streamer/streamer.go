package streamer

import (
	"context"
	"fmt"
	"github.com/Borislavv/video-streaming/internal/domain/dto"
	"github.com/Borislavv/video-streaming/internal/domain/entity"
	"github.com/Borislavv/video-streaming/internal/domain/logger"
	"github.com/Borislavv/video-streaming/internal/domain/repository"
	"github.com/Borislavv/video-streaming/internal/domain/service"
	"github.com/Borislavv/video-streaming/internal/domain/vo"
	"github.com/gorilla/websocket"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"gopkg.in/vansante/go-ffprobe.v2"
	"os"
	"strings"
	"sync"
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
	availableActionsMap = map[ActionEnum]struct{}{
		StreamByID:        {},
		StopStream:        {},
		DecreaseBufferCap: {},
	}
)

const (
	ProtoSeparator string = ":"

	StreamByID        ActionEnum = "ID"
	StopStream        ActionEnum = "stop"
	DecreaseBufferCap ActionEnum = "decrBuff"
)

type ResourceStreamer struct {
	ctx             context.Context
	logger          logger.Logger
	reader          service.Reader
	videoRepository repository.Video
}

func NewStreamingService(
	ctx context.Context,
	logger logger.Logger,
	reader service.Reader,
	videoRepository repository.Video,
) *ResourceStreamer {
	return &ResourceStreamer{
		ctx:             ctx,
		reader:          reader,
		logger:          logger,
		videoRepository: videoRepository,
	}
}

func (s *ResourceStreamer) Stream(conn *websocket.Conn) {
	s.logger.Info(fmt.Sprintf("[%v]: start streaming", conn.RemoteAddr()))

	wg := &sync.WaitGroup{}
	wg.Add(3)

	actionCh := make(chan Action, 1)
	decrBuffCh := make(chan struct{})

	go s.listenClient(wg, conn, actionCh)
	go s.handleBufferCapacity(wg, conn, decrBuffCh)
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
			do, data := s.parseMessage(b)
			if _, exists := availableActionsMap[do]; exists {
				actionsCh <- Action{do: do, data: data}
			} else {
				s.logger.Critical(fmt.Sprintf("do: %+v, data: %+v received unsupport action", do, data))
			}
		}
	}
}

func (s *ResourceStreamer) handleBufferCapacity(
	wg *sync.WaitGroup,
	conn *websocket.Conn,
	decrBuffCapCh <-chan struct{},
) {
	defer wg.Done()
	for range decrBuffCapCh {
		s.logger.Info(
			fmt.Sprintf(
				"[%v]: decreased buffer capacity action received",
				conn.RemoteAddr(),
			),
		)
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
			oid, e := primitive.ObjectIDFromHex(action.data)
			if e != nil {
				s.logger.Log(e)
				continue
			}
			v, e := s.videoRepository.Find(s.ctx, vo.ID{Value: oid})
			if e != nil {
				s.logger.Log(e)
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

func (s *ResourceStreamer) sendStartStreamMessage(resource entity.Resource, conn *websocket.Conn) error {
	audioCodec, videoCodec, err := s.codecs(resource)
	if err != nil {
		return s.logger.ErrorPropagate(fmt.Sprintf("[%v]: %v", conn.RemoteAddr(), err.Error()))
	}

	b := strings.Builder{}
	b.WriteString("start") // writing init. message identifier
	b.WriteString(ProtoSeparator)
	b.WriteString(audioCodec)
	b.WriteString(ProtoSeparator)
	b.WriteString(videoCodec)
	initMessage := b.String()

	// writing the stream initialization message in a websocket connection
	if err = conn.WriteMessage(websocket.TextMessage, []byte(initMessage)); err != nil {
		return s.logger.ErrorPropagate(fmt.Sprintf("[%v]: %v", conn.RemoteAddr(), err.Error()))
	}

	return nil
}

func (s *ResourceStreamer) sendStopStreamMessage(conn *websocket.Conn) error {
	if err := conn.WriteMessage(websocket.TextMessage, []byte(StopStream.String())); err != nil {
		return s.logger.CriticalPropagate(fmt.Sprintf("[%v]: %v", conn.RemoteAddr(), err.Error()))
	}
	return nil
}

func (s *ResourceStreamer) sendChunkStreamMessage(chunk dto.Chunk, conn *websocket.Conn) error {
	if chunk.GetError() != nil {
		return s.logger.CriticalPropagate(fmt.Sprintf("[%v]: %v", conn.RemoteAddr(), chunk.GetError().Error()))
	}

	if err := conn.WriteMessage(websocket.BinaryMessage, chunk.GetData()); err != nil {
		return s.logger.CriticalPropagate(fmt.Sprintf("[%v]: %v", conn.RemoteAddr(), err.Error()))
	}

	return nil
}

func (s *ResourceStreamer) streamResource(resource entity.Resource, conn *websocket.Conn) {
	if err := s.sendStartStreamMessage(resource, conn); err != nil {
		s.logger.Critical(fmt.Sprintf("[%v]: %v", conn.RemoteAddr(), err.Error()))
		return
	}

	for chunk := range s.reader.Read(resource) {
		if err := s.sendChunkStreamMessage(chunk, conn); err != nil {
			s.logger.Critical(fmt.Sprintf("[%v]: %v", conn.RemoteAddr(), err.Error()))
			break
			// TODO must be implemented method `sendErrorStreamMessage`
			// due to be able to tell client that error occurred on the server side
		}

		s.logger.Info(
			fmt.Sprintf(
				"[%v]: wrote %d bytes of '%v' to websocket",
				conn.RemoteAddr(), chunk.GetLen(), resource.Name,
			),
		)
	}

	if err := s.sendStopStreamMessage(conn); err != nil {
		s.logger.Critical(fmt.Sprintf("[%v]: %v", conn.RemoteAddr(), err.Error()))
		return
	}
}

func (s *ResourceStreamer) parseMessage(b []byte) (do ActionEnum, data string) {
	p := strings.Split(string(b), ProtoSeparator)
	if len(p) > 1 {
		return ActionEnum(p[0]), p[1]
	}
	return ActionEnum(p[0]), ""
}

// codecs will determine video and audio stream codecs of target resource
func (s *ResourceStreamer) codecs(
	resource entity.Resource,
) (
	audioCodec string,
	videoCodec string,
	e error,
) {
	file, err := os.Open(resource.GetFilepath())
	if err != nil {
		return "", "", s.logger.LogPropagate(err)
	}
	defer func() { _ = file.Close() }()

	data, err := ffprobe.ProbeReader(s.ctx, file)
	if err != nil {
		return "", "", s.logger.LogPropagate(err)
	}

	audioCodec = ""
	videoCodec = ""
	if data.FirstAudioStream() != nil {
		audioCodec = data.FirstAudioStream().CodecTagString
	}
	if data.FirstVideoStream() != nil {
		videoCodec = data.FirstVideoStream().CodecTagString
	}

	return audioCodec, videoCodec, nil
}
