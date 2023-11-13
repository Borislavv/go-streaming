package strategy

import (
	"context"
	"fmt"
	"github.com/Borislavv/video-streaming/internal/domain/dto"
	"github.com/Borislavv/video-streaming/internal/domain/entity"
	"github.com/Borislavv/video-streaming/internal/domain/errors"
	"github.com/Borislavv/video-streaming/internal/domain/logger"
	"github.com/Borislavv/video-streaming/internal/domain/repository"
	"github.com/Borislavv/video-streaming/internal/domain/service/tokenizer"
	"github.com/Borislavv/video-streaming/internal/domain/vo"
	"github.com/Borislavv/video-streaming/internal/infrastructure/service/detector"
	"github.com/Borislavv/video-streaming/internal/infrastructure/service/reader"
	"github.com/Borislavv/video-streaming/internal/infrastructure/service/streamer/action/enum"
	"github.com/Borislavv/video-streaming/internal/infrastructure/service/streamer/action/model"
	"github.com/Borislavv/video-streaming/internal/infrastructure/service/streamer/proto"
	"github.com/gorilla/websocket"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"math"
	"os"
)

type StreamByIDWithOffsetActionStrategy struct {
	ctx             context.Context
	logger          logger.Logger
	videoRepository repository.Video
	reader          reader.FileReader
	codecInfo       detector.Detector
	communicator    proto.Communicator
	tokenizer       tokenizer.Tokenizer
	chunkSize       int
}

func NewStreamByIDWithOffsetActionStrategy(
	ctx context.Context,
	logger logger.Logger,
	videoRepository repository.Video,
	reader reader.FileReader,
	codecInfo detector.Detector,
	communicator proto.Communicator,
	tokenizer tokenizer.Tokenizer,
	chunkSize int,
) *StreamByIDWithOffsetActionStrategy {
	return &StreamByIDWithOffsetActionStrategy{
		ctx:             ctx,
		logger:          logger,
		videoRepository: videoRepository,
		reader:          reader,
		codecInfo:       codecInfo,
		communicator:    communicator,
		tokenizer:       tokenizer,
		chunkSize:       chunkSize,
	}
}

// IsAppropriate - method will tell the service architect that the strategy is acceptable.
func (s *StreamByIDWithOffsetActionStrategy) IsAppropriate(action model.Action) bool {
	return action.Do == enum.StreamByIDWithOffset
}

// Do - will be streaming a target resource by ID from given offset.
func (s *StreamByIDWithOffsetActionStrategy) Do(action model.Action) error {
	// check the data is eligible
	data, ok := action.Data.(model.StreamByIdWithOffsetData)
	if !ok {
		return s.logger.CriticalPropagate(
			fmt.Errorf("'by id with offset' strategy cannot handle the given data '%+v'", data),
		)
	}

	// user authentication
	userID, err := s.tokenizer.Validate(data.Token)
	if err != nil {
		return s.logger.LogPropagate(err)
	}

	// parse the given video resource identifier
	oid, err := primitive.ObjectIDFromHex(data.ID)
	if err != nil {
		return s.logger.LogPropagate(err)
	}

	// searching the requested video resource
	v, err := s.videoRepository.FindOneByID(s.ctx, dto.NewVideoGetRequestDTO(vo.NewID(oid), userID))
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
	s.stream(v.Resource, data, action.Conn)

	return nil
}

func (s *StreamByIDWithOffsetActionStrategy) stream(
	resource entity.Resource,
	data model.StreamByIdWithOffsetData,
	conn *websocket.Conn,
) {
	audioCodec, videoCodec, err := s.codecInfo.Detect(resource)
	if err != nil {
		s.logger.Error(fmt.Sprintf("[%v]: %v", conn.RemoteAddr(), err.Error()))
		return
	}

	if err = s.communicator.Start(audioCodec, videoCodec, conn); err != nil {
		s.logger.Critical(fmt.Sprintf("[%v]: %v", conn.RemoteAddr(), err.Error()))
		return
	}

	file, err := os.Open(resource.GetFilepath())
	if err != nil {
		s.logger.Critical(fmt.Sprintf("[%v]: error resource opening: %v", conn.RemoteAddr(), err.Error()))
		return
	}
	defer func() { _ = file.Close() }()

	stat, err := file.Stat()
	if err != nil {
		s.logger.Critical(fmt.Sprintf("[%v]: error receiving resource stat: %v", conn.RemoteAddr(), err.Error()))
		return
	}

	totalChunks := stat.Size() / int64(s.chunkSize)

	chunkDuration := data.Duration / float64(totalChunks)

	targetChunk := math.Ceil(data.From / chunkDuration)

	offset := int64((s.chunkSize * int(targetChunk)) - s.chunkSize)

	for chunk := range s.reader.ReadByChunks(file, offset) {
		if err = s.communicator.Send(chunk, conn); err != nil {
			s.logger.Critical(fmt.Sprintf("[%v]: %v", conn.RemoteAddr(), err.Error()))
			break
		}

		s.logger.Info(
			fmt.Sprintf("[%v]: wrote %d bytes of '%v' to websocket",
				conn.RemoteAddr(), chunk.GetLen(), resource.Name,
			),
		)
	}

	if err = s.communicator.Stop(conn); err != nil {
		s.logger.Critical(fmt.Sprintf("[%v]: %v", conn.RemoteAddr(), err.Error()))
		return
	}
}
