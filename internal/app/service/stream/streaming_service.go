package stream

import (
	"context"
	"errors"
	"fmt"
	"github.com/Borislavv/video-streaming/internal/app/service"
	"github.com/Borislavv/video-streaming/internal/domain/model/video"
	"github.com/gorilla/websocket"
	"gopkg.in/vansante/go-ffprobe.v2"
	"log"
	"os"
	"sync"
)

const VideoPath = "/home/jared/jaredsplace/projects/go/streaming/internal/infrastructure/static/tmp/video/example_video_new.mp4"
const Video2Path = "/home/jared/jaredsplace/projects/go/streaming/internal/infrastructure/static/tmp/video/example_video_2_new.mp4"
const Video3Path = "/home/jared/jaredsplace/projects/go/streaming/internal/infrastructure/static/tmp/video/example_video_3_new.mp4"

type Action string

func (a Action) String() string {
	return string(a)
}

const (
	Start             Action = "start"
	Pause             Action = "pause"
	Stop              Action = "stop"
	Next              Action = "next"
	DecreaseBufferCap Action = "decrBuff"
)

type StreamingService struct {
	reader service.Reader
	errCh  chan error
}

func NewStreamingService(
	reader service.Reader,
	errCh chan error,
) *StreamingService {
	return &StreamingService{
		reader: reader,
		errCh:  errCh,
	}
}

func (s *StreamingService) Stream(conn *websocket.Conn) {
	log.Println("[streamer]: start streaming")

	wg := &sync.WaitGroup{}
	wg.Add(1)

	actionCh := make(chan Action, 1)
	decrBuffCapCh := make(chan struct{})

	go s.handleMessages(conn, actionCh, decrBuffCapCh)
	go s.handleBufferSize(decrBuffCapCh)
	go s.handleStream(wg, conn, actionCh)

	wg.Wait()

	log.Println("[streamer]: streaming is stopped")
}

func (s *StreamingService) handleStream(wg *sync.WaitGroup, conn *websocket.Conn, actionCh <-chan Action) {
	defer wg.Done()
	defer log.Println("handleStream: exit")

	videos := []*video.Video{
		video.New(VideoPath),
		video.New(Video2Path),
		video.New(Video3Path),
	}

	for action := range actionCh {
		if action == Next {
			l := len(videos)
			if l >= 1 {
				v := videos[0]
				if l == 1 {
					videos = []*video.Video{}
				} else {
					videos = append(videos[:0], videos[1:]...)
				}
				s.stream(v, conn)
			}
		}
	}
}

func (s *StreamingService) stream(video *video.Video, conn *websocket.Conn) {
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
		s.errCh <- err
		return
	}

	log.Println("Sent new start msg: " + startMsg)

	for chunk := range s.reader.Read(video) {
		if chunk.Err != nil {
			s.errCh <- chunk.Err
			continue
		}

		if err := conn.WriteMessage(websocket.BinaryMessage, chunk.Data); err != nil {
			s.errCh <- err
			continue
		}

		log.Printf("[streamer]: wrote %d bytes to websocket\n", chunk.Len)
	}

	if err := conn.WriteMessage(websocket.TextMessage, []byte(Stop.String())); err != nil {
		s.errCh <- err
		return
	}
}

func (s *StreamingService) handleBufferSize(decrBuffCapCh <-chan struct{}) {
	defer log.Println("handleBufferSize: exit")

	for range decrBuffCapCh {
		log.Println("----buffer must be decreased-----")
	}
}

func (s *StreamingService) handleMessages(
	conn *websocket.Conn,
	actionsCh chan<- Action,
	decrBuffCapCh chan<- struct{},
) {
	defer close(actionsCh)
	defer close(decrBuffCapCh)
	defer log.Println("handleMessages: exit")

	for {
		t, b, err := conn.ReadMessage()
		if err != nil {
			if websocket.IsCloseError(err, websocket.CloseGoingAway, websocket.CloseNormalClosure) {
				return
			}
			s.errCh <- err
			return
		}
		if t == websocket.TextMessage {
			action := Action(b)

			log.Printf("####ACTION### %s\n", action)

			if action == DecreaseBufferCap {
				decrBuffCapCh <- struct{}{}
				continue
			}
			if action == Start || action == Pause || action == Stop || action == Next {
				actionsCh <- action
				log.Printf("####ACTION### SEND %s\n", action)
				continue
			}
			s.errCh <- errors.New(fmt.Sprintf("found unknown action: %s", action))
		}
	}
}

func (s *StreamingService) codecs(video *video.Video) (a *ffprobe.Stream, v *ffprobe.Stream) {
	file, err := os.Open(video.GetPath())
	if err != nil {
		log.Fatalf("Failed to open file: %v", err)
	}
	defer file.Close()

	data, err := ffprobe.ProbeReader(context.Background(), file)
	if err != nil {
		log.Fatalln(err)
	}

	return data.FirstAudioStream(), data.FirstVideoStream()
}
