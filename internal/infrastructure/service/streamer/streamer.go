package streamer

import (
	"context"
	"fmt"
	"github.com/Borislavv/video-streaming/internal/domain/agg"
	"github.com/Borislavv/video-streaming/internal/domain/dto"
	"github.com/Borislavv/video-streaming/internal/domain/entity"
	"github.com/Borislavv/video-streaming/internal/domain/logger"
	"github.com/Borislavv/video-streaming/internal/domain/repository"
	"github.com/Borislavv/video-streaming/internal/domain/service"
	"github.com/gorilla/websocket"
	"gopkg.in/vansante/go-ffprobe.v2"
	"os"
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
	wg.Add(1)

	actionCh := make(chan Action, 1)
	decrBuffCapCh := make(chan struct{})

	go s.listenClient(conn, actionCh, decrBuffCapCh)
	go s.handleBufferSize(decrBuffCapCh)
	go s.handleActions(wg, conn, actionCh)

	wg.Wait()

	s.logger.Info("streaming is stopped")
}

func (s *ResourceStreamer) handleActions(wg *sync.WaitGroup, conn *websocket.Conn, actionCh <-chan Action) {
	defer wg.Done()

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
		if action == Start {
			s.logger.Info("action 'start' received")
			v := videos[c]
			s.logger.Info(fmt.Sprintf("streaming 'video':'%v'", v.Name))
			s.stream(v, conn)
		} else if action == Next {
			if c < l {
				c++
			}
			s.logger.Info("action 'next' received")
			v := videos[c]
			s.logger.Info(fmt.Sprintf("streaming 'video':'%v'", v.Name))
			s.stream(v, conn)
		} else if action == Previous {
			if c >= 1 {
				c--
			}
			s.logger.Info("action 'previous' received")
			v := videos[c]
			s.logger.Info(fmt.Sprintf("streaming 'video':'%v'", v.Name))
			s.stream(v, conn)
		}
	}
}

func (s *ResourceStreamer) stream(video *agg.Video, conn *websocket.Conn) {
	startMsg := Start.String()

	ac, vc := s.codecs(video.Resource)
	if ac != nil {
		startMsg += ":" + ac.CodecTagString
	} else {
		startMsg += ":"
	}
	if vc != nil {
		startMsg += ":" + vc.CodecTagString
	} else {
		startMsg += ":"
	}

	if err := conn.WriteMessage(websocket.TextMessage, []byte(startMsg)); err != nil {
		s.logger.Error(err)
		return
	}

	for chunk := range s.reader.Read(video.Resource) {
		if chunk.Err != nil {
			s.logger.Critical(chunk.Err)
			continue
		}

		if err := conn.WriteMessage(websocket.BinaryMessage, chunk.Data); err != nil {
			s.logger.Error(err)
			continue
		}

		s.logger.Info(fmt.Sprintf("wrote %d bytes of '%v' to websocket", chunk.Len, video.Name))
	}

	if err := conn.WriteMessage(websocket.TextMessage, []byte(Stop.String())); err != nil {
		s.logger.Error(err)
		return
	}
}

func (s *ResourceStreamer) handleBufferSize(decrBuffCapCh <-chan struct{}) {
	for range decrBuffCapCh {
		s.logger.Info("decreased buffer capacity action received")
	}
}

func (s *ResourceStreamer) listenClient(
	conn *websocket.Conn,
	actionsCh chan<- Action,
	decrBuffCapCh chan<- struct{},
) {
	defer close(actionsCh)
	defer close(decrBuffCapCh)

	for {
		t, b, err := conn.ReadMessage()
		if err != nil {
			if websocket.IsCloseError(err, websocket.CloseGoingAway, websocket.CloseNormalClosure) {
				return
			}
			s.logger.Error(err)
			return
		}
		if t == websocket.TextMessage {
			action := Action(b)

			if action == DecreaseBufferCap {
				decrBuffCapCh <- struct{}{}
				continue
			}
			if action == Start || action == Stop || action == Next || action == Previous {
				actionsCh <- action
				continue
			}
			s.logger.Emergency(fmt.Sprintf("found unknown action: %s", action))
		}
	}
}

func (s *ResourceStreamer) codecs(resource entity.Resource) (audioStream *ffprobe.Stream, videoStream *ffprobe.Stream) {
	file, err := os.Open(resource.GetFilepath())
	if err != nil {
		s.logger.Emergency(err)
	}
	defer func() {
		if err = file.Close(); err != nil {
			s.logger.Emergency(err)
		}
	}()

	data, err := ffprobe.ProbeReader(context.Background(), file)
	if err != nil {
		s.logger.Emergency(err)
	}

	return data.FirstAudioStream(), data.FirstVideoStream()
}
