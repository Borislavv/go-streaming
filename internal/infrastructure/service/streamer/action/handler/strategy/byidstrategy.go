package strategy

import (
	"context"
	"fmt"
	"github.com/Borislavv/video-streaming/internal/domain/entity"
	"github.com/Borislavv/video-streaming/internal/domain/errors"
	"github.com/Borislavv/video-streaming/internal/domain/logger"
	"github.com/Borislavv/video-streaming/internal/domain/repository"
	"github.com/Borislavv/video-streaming/internal/domain/vo"
	"github.com/Borislavv/video-streaming/internal/infrastructure/service/codec"
	"github.com/Borislavv/video-streaming/internal/infrastructure/service/reader"
	"github.com/Borislavv/video-streaming/internal/infrastructure/service/streamer/action/enum"
	"github.com/Borislavv/video-streaming/internal/infrastructure/service/streamer/action/model"
	"github.com/Borislavv/video-streaming/internal/infrastructure/service/streamer/proto"
	"github.com/gorilla/websocket"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type StreamByIDActionStrategy struct {
	ctx             context.Context
	logger          logger.Logger
	videoRepository repository.Video
	reader          reader.Reader
	codecInfo       codec.Detector
	communicator    proto.Communicator
}

func NewStreamByIDActionStrategy(
	ctx context.Context,
	logger logger.Logger,
	videoRepository repository.Video,
	reader reader.Reader,
	codecInfo codec.Detector,
	communicator proto.Communicator,
) *StreamByIDActionStrategy {
	return &StreamByIDActionStrategy{
		ctx:             ctx,
		logger:          logger,
		videoRepository: videoRepository,
		reader:          reader,
		codecInfo:       codecInfo,
		communicator:    communicator,
	}
}

func (s *StreamByIDActionStrategy) IsAppropriate(action model.Action) bool {
	return action.Do == enum.StreamByID
}

func (s *StreamByIDActionStrategy) Do(action model.Action) error {
	oid, err := primitive.ObjectIDFromHex(action.Data)
	if err != nil {
		return s.logger.LogPropagate(err)
	}
	v, err := s.videoRepository.Find(s.ctx, vo.ID{Value: oid})
	if err != nil {
		if errors.IsEntityNotFoundError(err) {
			if err = s.communicator.Error(err, action.Conn); err != nil {
				return s.logger.LogPropagate(err)
			}
		}
		return s.logger.LogPropagate(err)
	}
	s.logger.Info(fmt.Sprintf("[%v]: streaming 'resource':'%v'", action.Conn.RemoteAddr(), v.Resource.Name))

	s.stream(v.Resource, action.Conn)

	return nil
}

// todo need to think about ctx for stream for to be able stop it and skip current action
func (s *StreamByIDActionStrategy) stream(resource entity.Resource, conn *websocket.Conn) {
	audioCodec, videoCodec, err := s.codecInfo.Detect(resource)
	if err != nil {
		s.logger.Error(fmt.Sprintf("[%v]: %v", conn.RemoteAddr(), err.Error()))
		return
	}

	if err = s.communicator.Start(audioCodec, videoCodec, conn); err != nil {
		s.logger.Critical(fmt.Sprintf("[%v]: %v", conn.RemoteAddr(), err.Error()))
		return
	}

	for chunk := range s.reader.Read(resource) {
		if err = s.communicator.Send(chunk, conn); err != nil {
			s.logger.Critical(fmt.Sprintf("[%v]: %v", conn.RemoteAddr(), err.Error()))
			break
		}

		s.logger.Info(fmt.Sprintf("[%v]: wrote %d bytes of '%v' to websocket",
			conn.RemoteAddr(), chunk.GetLen(), resource.Name,
		))
	}

	if err = s.communicator.Stop(conn); err != nil {
		s.logger.Critical(fmt.Sprintf("[%v]: %v", conn.RemoteAddr(), err.Error()))
		return
	}
}
