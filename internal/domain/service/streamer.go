package service

import (
	"context"
	"fmt"
	"github.com/Borislavv/video-streaming/internal/domain/entity"
	"github.com/gorilla/websocket"
	"gopkg.in/vansante/go-ffprobe.v2"
	"os"
	"sync"
)

const VideoPath = "/streaming_root/internal/infrastructure/static/tmp/video/example_video_new.mp4"
const Video2Path = "/streaming_root/internal/infrastructure/static/tmp/video/example_video_2_new.mp4"
const Video3Path = "/streaming_root/internal/infrastructure/static/tmp/video/example_video_3_new.mp4"

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
	reader Reader
	logger Logger
}

func NewStreamingService(
	reader Reader,
	logger Logger,
) *ResourceStreamer {
	return &ResourceStreamer{
		reader: reader,
		logger: logger,
	}
}

func (s *ResourceStreamer) Stream(conn *websocket.Conn) {
	s.logger.Info("[streamer]: start streaming")

	wg := &sync.WaitGroup{}
	wg.Add(1)

	actionCh := make(chan Action, 1)
	decrBuffCapCh := make(chan struct{})

	go s.handleMessages(conn, actionCh, decrBuffCapCh)
	go s.handleBufferSize(decrBuffCapCh)
	go s.handleStream(wg, conn, actionCh)

	wg.Wait()

	s.logger.Info("[streamer]: streaming is stopped")
}

func (s *ResourceStreamer) handleStream(wg *sync.WaitGroup, conn *websocket.Conn, actionCh <-chan Action) {
	defer wg.Done()
	defer s.logger.Info("handleStream: exit")

	videos := []*entity.Video{
		1: entity.NewVideo("1", VideoPath, "1"),
		2: entity.NewVideo("2", Video2Path, "2"),
		3: entity.NewVideo("3", Video3Path, "3"),
	}

	l := len(videos) - 1
	c := 0
	for action := range actionCh {
		if action == Next {
			if c < l {
				c++
			}

			v := videos[c]
			s.stream(v, conn)
		} else if action == Previous {
			if c > 1 {
				c--
			}

			v := videos[c]
			s.stream(v, conn)
		}
	}
}

func (s *ResourceStreamer) stream(video *entity.Video, conn *websocket.Conn) {
	startMsg := Start.String()

	ac, vc := s.codecs(video)
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

	for chunk := range s.reader.Read(video) {
		if chunk.Err != nil {
			s.logger.Critical(chunk.Err)
			continue
		}

		if err := conn.WriteMessage(websocket.BinaryMessage, chunk.Data); err != nil {
			s.logger.Error(err)
			continue
		}

		s.logger.Info(fmt.Sprintf("[streamer]: wrote %d bytes to websocket", chunk.Len))
	}

	if err := conn.WriteMessage(websocket.TextMessage, []byte(Stop.String())); err != nil {
		s.logger.Error(err)
		return
	}
}

func (s *ResourceStreamer) handleBufferSize(decrBuffCapCh <-chan struct{}) {
	defer s.logger.Info("handleBufferSize: exit")

	for range decrBuffCapCh {
		s.logger.Info("Decreased buffer capacity")
	}
}

func (s *ResourceStreamer) handleMessages(
	conn *websocket.Conn,
	actionsCh chan<- Action,
	decrBuffCapCh chan<- struct{},
) {
	defer close(actionsCh)
	defer close(decrBuffCapCh)
	defer s.logger.Info("handleMessages: exit")

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

func (s *ResourceStreamer) codecs(video *entity.Video) (a *ffprobe.Stream, v *ffprobe.Stream) {
	file, err := os.Open(video.GetPath())
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
