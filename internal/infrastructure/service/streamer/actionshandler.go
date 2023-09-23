package streamer

import (
	"context"
	"fmt"
	"github.com/Borislavv/video-streaming/internal/domain/entity"
	"github.com/Borislavv/video-streaming/internal/domain/errs"
	"github.com/Borislavv/video-streaming/internal/domain/logger"
	"github.com/Borislavv/video-streaming/internal/domain/repository"
	"github.com/Borislavv/video-streaming/internal/domain/service"
	"github.com/Borislavv/video-streaming/internal/domain/vo"
	"github.com/gorilla/websocket"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"sync"
)

type WebSocketActionsHandler struct {
	ctx             context.Context
	logger          logger.Logger
	reader          service.Reader
	videoRepository repository.Video
	communicator    ProtoCommunicator
	codecInfo       CodecsDetector
}

func NewWebSocketActionsHandler(
	ctx context.Context,
	logger logger.Logger,
	reader service.Reader,
	videoRepository repository.Video,
	communicator ProtoCommunicator,
	codecInfo CodecsDetector,
) *WebSocketActionsHandler {
	return &WebSocketActionsHandler{
		ctx:             ctx,
		logger:          logger,
		reader:          reader,
		videoRepository: videoRepository,
		communicator:    communicator,
		codecInfo:       codecInfo,
	}
}

func (h *WebSocketActionsHandler) Handle(wg *sync.WaitGroup, conn *websocket.Conn, actionsCh <-chan Action) {
	wg.Add(1)
	go func() {
		defer wg.Done()

		for action := range actionsCh {
			// todo must be moved in the strategies of actions and in the 'stream' action must be 'stream' strategies
			switch action.do {
			case streamByID:
				oid, err := primitive.ObjectIDFromHex(action.data)
				if err != nil {
					h.logger.Log(err)
					continue
				}
				v, err := h.videoRepository.Find(h.ctx, vo.ID{Value: oid})
				if err != nil {
					h.logger.Log(err)
					if errs.IsNotFoundError(err) {
						if err = h.communicator.Error(err, conn); err != nil {
							h.logger.Log(err)
						}
						continue
					}
					continue
				}
				h.logger.Info(fmt.Sprintf("[%v]: streaming 'resource':'%v'", conn.RemoteAddr(), v.Resource.Name))

				// todo need to think about ctx for stream for to be able stop it and skip current action
				h.stream(v.Resource, conn)
			}
		}
	}()
}

func (h *WebSocketActionsHandler) stream(resource entity.Resource, conn *websocket.Conn) {
	audioCodec, videoCodec, err := h.codecInfo.Detect(resource)
	if err != nil {
		h.logger.Error(fmt.Sprintf("[%v]: %v", conn.RemoteAddr(), err.Error()))
		return
	}

	if err = h.communicator.Start(audioCodec, videoCodec, conn); err != nil {
		h.logger.Critical(fmt.Sprintf("[%v]: %v", conn.RemoteAddr(), err.Error()))
		return
	}

	for chunk := range h.reader.Read(resource) {
		if err = h.communicator.Send(chunk, conn); err != nil {
			h.logger.Critical(fmt.Sprintf("[%v]: %v", conn.RemoteAddr(), err.Error()))
			break
		}

		h.logger.Info(fmt.Sprintf("[%v]: wrote %d bytes of '%v' to websocket",
			conn.RemoteAddr(), chunk.GetLen(), resource.Name,
		))
	}

	if err = h.communicator.Stop(conn); err != nil {
		h.logger.Critical(fmt.Sprintf("[%v]: %v", conn.RemoteAddr(), err.Error()))
		return
	}
}
