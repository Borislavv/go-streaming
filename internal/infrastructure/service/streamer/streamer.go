package streamer

import (
	"context"
	"fmt"
	"github.com/Borislavv/video-streaming/internal/domain/dto"
	"github.com/Borislavv/video-streaming/internal/domain/entity"
	"github.com/Borislavv/video-streaming/internal/domain/logger"
	"github.com/Borislavv/video-streaming/internal/domain/repository"
	"github.com/Borislavv/video-streaming/internal/domain/service"
	"github.com/gorilla/websocket"
	"gopkg.in/vansante/go-ffprobe.v2"
	"os"
	"strings"
	"sync"
)

type Action string

func (a Action) String() string {
	return string(a)
}

const (
	Start             Action = "start"
	Stop              Action = "stop"
	Next              Action = "next"
	Previous          Action = "prev"
	DecreaseBufferCap Action = "decrBuff"
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
	s.logger.Info("start streaming")

	wg := &sync.WaitGroup{}
	wg.Add(3)

	actionCh := make(chan Action, 1)
	decrBuffCh := make(chan struct{})

	go s.listenClient(wg, conn, actionCh)
	go s.handleBufferSize(wg, decrBuffCh)
	go s.handleActions(wg, conn, actionCh, decrBuffCh)

	wg.Wait()

	s.logger.Info("streaming is stopped")
}

func (s *ResourceStreamer) handleActions(
	wg *sync.WaitGroup,
	conn *websocket.Conn,
	actionCh <-chan Action,
	decrBuffCapCh chan<- struct{},
) {
	defer func() {
		close(decrBuffCapCh)
		wg.Done()
	}()

	videos, err := s.videoRepository.FindList(
		s.ctx,
		&dto.VideoListRequestDto{
			PaginationRequestDto: dto.PaginationRequestDto{
				Page:  1,
				Limit: 10,
			},
		},
	)
	if err != nil {
		s.logger.Log(err)
		return
	}

	l := len(videos) - 1
	c := 0
	for action := range actionCh {
		switch action {
		case Start:
			s.logger.Info("action 'start' received")
		case Next:
			if c < l {
				c++
			}
			s.logger.Info("action 'next' received")
		case Previous:
			if c >= 1 {
				c--
			}
			s.logger.Info("action 'previous' received")
		case DecreaseBufferCap:
			decrBuffCapCh <- struct{}{}
			continue
		}

		resource := videos[c].Resource
		s.logger.Info(fmt.Sprintf("streaming 'resource':'%v'", resource.Name))
		s.stream(resource, conn)
	}
}

func (s *ResourceStreamer) sendStartStreamMessage(resource entity.Resource, conn *websocket.Conn) error {
	audioCodec, videoCodec, err := s.codecs(resource)
	if err != nil {
		return s.logger.LogPropagate(err)
	}

	b := strings.Builder{}
	b.WriteString(Start.String()) // writing init. message identifier
	b.WriteString(":")
	b.WriteString(audioCodec)
	b.WriteString(":")
	b.WriteString(videoCodec)
	initMessage := b.String()

	// writing the stream initialization message in a websocket connection
	if err = conn.WriteMessage(websocket.TextMessage, []byte(initMessage)); err != nil {
		return s.logger.LogPropagate(err)
	}

	return nil
}

func (s *ResourceStreamer) sendStopStreamMessage(conn *websocket.Conn) error {
	if err := conn.WriteMessage(websocket.TextMessage, []byte(Stop.String())); err != nil {
		return s.logger.CriticalPropagate(err)
	}
	return nil
}

func (s *ResourceStreamer) sendChunkStreamMessage(chunk *dto.Chunk, conn *websocket.Conn) error {
	if chunk.Err != nil {
		return s.logger.CriticalPropagate(chunk.Err)
	}

	if err := conn.WriteMessage(websocket.BinaryMessage, chunk.Data); err != nil {
		return s.logger.CriticalPropagate(err)
	}

	return nil
}

func (s *ResourceStreamer) stream(resource entity.Resource, conn *websocket.Conn) {
	if err := s.sendStartStreamMessage(resource, conn); err != nil {
		s.logger.Critical(err)
		return
	}

	for chunk := range s.reader.Read(resource) {
		if err := s.sendChunkStreamMessage(chunk, conn); err != nil {
			s.logger.Critical(err)
			break
			// TODO must be implemented method `sendErrorStreamMessage` due to be able to tell client that error occurred on the server side
		}
		s.logger.Info(fmt.Sprintf("wrote %d bytes of '%v' to websocket", chunk.Len, resource.Name))
	}

	if err := s.sendStopStreamMessage(conn); err != nil {
		s.logger.Critical(err)
		return
	}
}

func (s *ResourceStreamer) handleBufferSize(wg *sync.WaitGroup, decrBuffCapCh <-chan struct{}) {
	defer wg.Done()
	for range decrBuffCapCh {
		s.logger.Info("decreased buffer capacity action received")
	}
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
				s.logger.Info("websocket connection has been closed")
				return
			}
			s.logger.Error(err)
			return
		}
		if t == websocket.TextMessage {
			actionsCh <- Action(b)
		}
	}
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
