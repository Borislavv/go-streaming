package strategy

import (
	"context"
	"fmt"
	"github.com/Borislavv/video-streaming/internal/domain/dto"
	"github.com/Borislavv/video-streaming/internal/domain/entity"
	"github.com/Borislavv/video-streaming/internal/domain/errors"
	"github.com/Borislavv/video-streaming/internal/domain/logger/interface"
	repositoryinterface "github.com/Borislavv/video-streaming/internal/domain/repository/interface"
	"github.com/Borislavv/video-streaming/internal/domain/service/di/interface"
	tokenizerinterface "github.com/Borislavv/video-streaming/internal/domain/service/tokenizer/interface"
	"github.com/Borislavv/video-streaming/internal/domain/vo"
	detectorinterface "github.com/Borislavv/video-streaming/internal/infrastructure/service/detector/interface"
	readerinterface "github.com/Borislavv/video-streaming/internal/infrastructure/service/reader/interface"
	"github.com/Borislavv/video-streaming/internal/infrastructure/service/streamer/action/enum"
	"github.com/Borislavv/video-streaming/internal/infrastructure/service/streamer/action/model"
	protointerface "github.com/Borislavv/video-streaming/internal/infrastructure/service/streamer/proto/interface"
	"github.com/gorilla/websocket"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"os"
)

const zeroOffset = 0

type StreamByIDActionStrategy struct {
	ctx             context.Context
	logger          loggerinterface.Logger
	videoRepository repositoryinterface.Video
	reader          readerinterface.FileReader
	codecInfo       detectorinterface.Codecs
	communicator    protointerface.Communicator
	tokenizer       tokenizerinterface.Tokenizer
}

func NewStreamByIDActionStrategy(serviceContainer diinterface.ContainerManager) (*StreamByIDActionStrategy, error) {
	loggerService, err := serviceContainer.GetLoggerService()
	if err != nil {
		return nil, err
	}

	ctx, err := serviceContainer.GetCtx()
	if err != nil {
		return nil, loggerService.LogPropagate(err)
	}

	videoRepository, err := serviceContainer.GetVideoRepository()
	if err != nil {
		return nil, loggerService.LogPropagate(err)
	}

	fileReader, err := serviceContainer.GetFileReaderService()
	if err != nil {
		return nil, loggerService.LogPropagate(err)
	}

	codecsDetector, err := serviceContainer.GetCodecsDetectorService()
	if err != nil {
		return nil, loggerService.LogPropagate(err)
	}

	webSocketCommunicator, err := serviceContainer.GetWebSocketCommunicatorService()
	if err != nil {
		return nil, loggerService.LogPropagate(err)
	}

	tokenizerService, err := serviceContainer.GetTokenizerService()
	if err != nil {
		return nil, loggerService.LogPropagate(err)
	}

	return &StreamByIDActionStrategy{
		ctx:             ctx,
		logger:          loggerService,
		videoRepository: videoRepository,
		reader:          fileReader,
		codecInfo:       codecsDetector,
		communicator:    webSocketCommunicator,
		tokenizer:       tokenizerService,
	}, nil
}

// IsAppropriate - method will tell the service architect that the strategy is acceptable.
func (s *StreamByIDActionStrategy) IsAppropriate(action model.Action) bool {
	return action.Do == enum.StreamByID
}

// Do - will be streaming a target resource by ID.
func (s *StreamByIDActionStrategy) Do(action model.Action) error {
	// check the data is eligible
	data, ok := action.Data.(*model.StreamByIdData)
	if !ok {
		return s.logger.CriticalPropagate(
			fmt.Errorf("'by id' strategy cannot handle the given data '%+v'", data),
		)
	}

	// user authentication
	userID, err := s.tokenizer.Verify(data.Token)
	if err != nil {
		return s.logger.LogPropagate(err)
	}

	// parse the given video resource identifier
	oid, err := primitive.ObjectIDFromHex(data.ID)
	if err != nil {
		return s.logger.LogPropagate(err)
	}

	// find the target resource
	q := dto.NewVideoGetRequestDTO(vo.NewID(oid), "", vo.ID{}, userID)
	v, err := s.videoRepository.FindOneByID(s.ctx, q)
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
