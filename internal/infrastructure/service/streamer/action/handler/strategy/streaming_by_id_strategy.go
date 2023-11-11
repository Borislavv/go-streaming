package strategy

import (
	"context"
	"fmt"
	"github.com/Borislavv/video-streaming/internal/domain/dto"
	"github.com/Borislavv/video-streaming/internal/domain/entity"
	"github.com/Borislavv/video-streaming/internal/domain/errors"
	"github.com/Borislavv/video-streaming/internal/domain/logger"
	repository "github.com/Borislavv/video-streaming/internal/domain/repository/storage"
	"github.com/Borislavv/video-streaming/internal/domain/vo"
	"github.com/Borislavv/video-streaming/internal/infrastructure/service/detector"
	"github.com/Borislavv/video-streaming/internal/infrastructure/service/reader"
	"github.com/Borislavv/video-streaming/internal/infrastructure/service/streamer/action/enum"
	"github.com/Borislavv/video-streaming/internal/infrastructure/service/streamer/action/model"
	"github.com/Borislavv/video-streaming/internal/infrastructure/service/streamer/proto"
	"github.com/gorilla/websocket"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"os"
)

const zeroOffset = 0

type StreamByIDActionStrategy struct {
	ctx             context.Context
	logger          logger.Logger
	videoRepository repository.Video
	reader          reader.FileReader
	codecInfo       detector.Detector
	communicator    proto.Communicator
}

func NewStreamByIDActionStrategy(
	ctx context.Context,
	logger logger.Logger,
	videoRepository repository.Video,
	reader reader.FileReader,
	codecInfo detector.Detector,
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

// IsAppropriate - method will tell the service architect that the strategy is acceptable.
func (s *StreamByIDActionStrategy) IsAppropriate(action model.Action) bool {
	return action.Do == enum.StreamByID
}

// Do - will be streaming a target resource by ID.
func (s *StreamByIDActionStrategy) Do(action model.Action) error {
	// check the data is eligible
	data, ok := action.Data.(model.StreamByIdData)
	if !ok {
		return s.logger.CriticalPropagate(
			fmt.Errorf("'by id' strategy cannot handle the given data '%+v'", data),
		)
	}

	// parse the given video resource identifier
	oid, err := primitive.ObjectIDFromHex(data.ID)
	if err != nil {
		return s.logger.LogPropagate(err)
	}

	// TODO need to past token from the client and parse it on each request (actually fetch from cache)
	// find the target resource
	v, err := s.videoRepository.FindOneByID(s.ctx, dto.NewVideoGetRequestDTO(vo.ID{Value: oid}, vo.ID{Value: oid})) // TODO this will not work at now because video oid passed as user oid
	if err != nil {
		if errors.IsEntityNotFoundError(err) {
			if err = s.communicator.Error(err, action.Conn); err != nil {
				return s.logger.LogPropagate(err)
			}
		}
		return s.logger.LogPropagate(err)
	}
	s.logger.Info(fmt.Sprintf("[%v]: streaming 'resource':'%v'", action.Conn.RemoteAddr(), v.Resource.Name))

	// video resource streaming
	s.stream(v.Resource, action.Conn)

	return nil
}

// stream - the method which composed all useful work of really streaming.
func (s *StreamByIDActionStrategy) stream(resource entity.Resource, conn *websocket.Conn) {
	// detect the audio and video codecs
	audioCodec, videoCodec, err := s.codecInfo.Detect(resource)
	if err != nil {
		s.logger.Error(fmt.Sprintf("[%v]: %v", conn.RemoteAddr(), err.Error()))
		return
	}

	// send the initializing message to client side
	if err = s.communicator.Start(audioCodec, videoCodec, conn); err != nil {
		s.logger.Critical(fmt.Sprintf("[%v]: %v", conn.RemoteAddr(), err.Error()))
		return
	}

	// open the target resource file
	file, err := os.Open(resource.GetFilepath())
	if err != nil {
		s.logger.Critical(fmt.Sprintf("[%v]: error resource opening: %v", conn.RemoteAddr(), err.Error()))
		return
	}
	defer func() { _ = file.Close() }()

	// read the whole target file
	//chunk := s.reader.ReadAll(file)
	//// send the received chunk which is contains whole file
	//if err = s.communicator.Send(chunk, conn); err != nil {
	//	s.logger.Critical(fmt.Sprintf("[%v]: %v", conn.RemoteAddr(), err))
	//	return
	//}
	//s.logger.Info(
	//	fmt.Sprintf("[%v]: wrote %d bytes of '%v' to websocket",
	//		conn.RemoteAddr(), chunk.GetLen(), resource.Name,
	//	),
	//)

	// read the target file by chunks from zero offset
	for chunk := range s.reader.ReadByChunks(file, zeroOffset) {
		if err = s.communicator.Send(chunk, conn); err != nil {
			s.logger.Critical(fmt.Sprintf("[%v]: %v", conn.RemoteAddr(), err))
			break
		}

		s.logger.Info(
			fmt.Sprintf("[%v]: wrote %d bytes of '%v' to websocket",
				conn.RemoteAddr(), chunk.GetLen(), resource.Name,
			),
		)
	}

	// stop the streaming by sending appropriate message to client side
	if err = s.communicator.Stop(conn); err != nil {
		s.logger.Critical(fmt.Sprintf("[%v]: %v", conn.RemoteAddr(), err.Error()))
		return
	}
}
