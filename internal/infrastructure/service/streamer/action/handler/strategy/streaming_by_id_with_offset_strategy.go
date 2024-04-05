package strategy

import (
	"context"
	"fmt"
	"github.com/Borislavv/video-streaming/internal/domain/dto"
	"github.com/Borislavv/video-streaming/internal/domain/entity"
	"github.com/Borislavv/video-streaming/internal/domain/errtype"
	"github.com/Borislavv/video-streaming/internal/domain/logger/interface"
	repositoryinterface "github.com/Borislavv/video-streaming/internal/domain/repository/interface"
	tokenizerinterface "github.com/Borislavv/video-streaming/internal/domain/service/tokenizer/interface"
	"github.com/Borislavv/video-streaming/internal/domain/vo"
	detectorinterface "github.com/Borislavv/video-streaming/internal/infrastructure/service/detector/interface"
	readerinterface "github.com/Borislavv/video-streaming/internal/infrastructure/service/reader/interface"
	"github.com/Borislavv/video-streaming/internal/infrastructure/service/streamer/action/enum"
	"github.com/Borislavv/video-streaming/internal/infrastructure/service/streamer/action/model"
	protointerface "github.com/Borislavv/video-streaming/internal/infrastructure/service/streamer/proto/interface"
	"github.com/gorilla/websocket"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"math"
	"os"
)

type StreamByIDWithOffsetActionStrategy struct {
	ctx             context.Context
	logger          loggerinterface.Logger
	videoRepository repositoryinterface.Video
	reader          readerinterface.FileReader
	codecInfo       detectorinterface.Codecs
	communicator    protointerface.Communicator
	tokenizer       tokenizerinterface.Tokenizer
	chunkSize       int
}

func NewStreamByIDWithOffsetActionStrategy(
	ctx context.Context,
	logger loggerinterface.Logger,
	videoRepository repositoryinterface.Video,
	reader readerinterface.FileReader,
	codecInfo detectorinterface.Codecs,
	communicator protointerface.Communicator,
	tokenizer tokenizerinterface.Tokenizer,
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
	data, ok := action.Data.(*model.StreamByIdWithOffsetData)
	if !ok {
		return s.logger.CriticalPropagate(
			fmt.Errorf("'by id with offset' strategy cannot handle the given data '%+v'", data),
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

	// searching the requested video resource
	q := dto.NewVideoGetRequestDTO(vo.NewID(oid), "", vo.ID{}, userID)
	v, err := s.videoRepository.FindOneByID(s.ctx, q)
	if err != nil {
		if errtype.IsEntityNotFoundError(err) {
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
	data *model.StreamByIdWithOffsetData,
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
