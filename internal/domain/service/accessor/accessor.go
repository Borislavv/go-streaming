package accessor

import (
	"context"
	"fmt"
	"github.com/Borislavv/video-streaming/internal/domain/agg"
	"github.com/Borislavv/video-streaming/internal/domain/errors"
	"github.com/Borislavv/video-streaming/internal/domain/logger"
	"github.com/Borislavv/video-streaming/internal/domain/repository"
	"github.com/Borislavv/video-streaming/internal/domain/vo"
	"reflect"
)

type AggregateAccessType int
type AggregateAccessHandler func(userID vo.ID, agg agg.Aggregate) (isGranted bool, err error)
type AggregateAccessIsAppropriateHandler func(agg agg.Aggregate) (isSupported bool)

const (
	VideoType AggregateAccessType = iota
	AudioType
	ResourceType
	UserType
)

type AccessService struct {
	ctx                       context.Context
	logger                    logger.Logger
	videoRepository           repository.Video
	audioRepository           repository.Audio
	userRepository            repository.User
	resourceRepository        repository.Resource
	handlers                  map[AggregateAccessType]AggregateAccessHandler
	isAppropriateHandlerFuncs map[AggregateAccessType]AggregateAccessIsAppropriateHandler
}

func NewAccessService(
	ctx context.Context,
	logger logger.Logger,
	videoRepository repository.Video,
	audioRepository repository.Audio,
	userRepository repository.User,
	resourceRepository repository.Resource,
) *AccessService {
	return (&AccessService{
		ctx:                       ctx,
		logger:                    logger,
		videoRepository:           videoRepository,
		audioRepository:           audioRepository,
		userRepository:            userRepository,
		resourceRepository:        resourceRepository,
		handlers:                  map[AggregateAccessType]AggregateAccessHandler{},
		isAppropriateHandlerFuncs: map[AggregateAccessType]AggregateAccessIsAppropriateHandler{},
	}).setHandlers()
}

// IsGranted is a method which will check the access to target scope of aggregates.
func (s *AccessService) IsGranted(userID vo.ID, aggregates []agg.Aggregate) (isGranted bool, err error) {
	for _, aggregate := range aggregates {
		for aggregateAccessType, isAppropriateHandler := range s.isAppropriateHandlerFuncs {
			if isAppropriateHandler(aggregate) {
				appropriateHandler, found := s.handlers[aggregateAccessType]
				if !found {
					return false, s.logger.LogPropagate(
						fmt.Errorf(
							"appropriate handler was not found by type '%d' into the handlers map",
							aggregateAccessType,
						),
					)
				}

				isGranted, err = appropriateHandler(userID, aggregate)
				if err != nil {
					return false, s.logger.LogPropagate(err)
				}
				if !isGranted {
					return false, nil
				}
			}
		}
	}

	return true, nil
}

// video
func (s *AccessService) videoHandler(userID vo.ID, aggregate agg.Aggregate) (isGranted bool, err error) {
	videoAgg, ok := aggregate.(*agg.Video)
	if !ok {
		return false, s.logger.LogPropagate(
			fmt.Errorf(
				"unable to check access for given aggregate of type '%v' in video access handler",
				reflect.TypeOf(aggregate).Name(),
			),
		)
	}

	userAgg, err := s.userRepository.Find(s.ctx, userID)
	if err != nil {
		return false, s.logger.LogPropagate(err)
	}

	for _, videoID := range userAgg.VideoIDs {
		if videoID.Value == videoAgg.ID.Value {
			return true, nil
		}
	}

	return true, nil
}
func (s *AccessService) audioIsAppropriateHandler(aggregate agg.Aggregate) (isAppropriate bool) {
	if _, ok := aggregate.(*agg.Audio); ok {
		return true
	}
	return false
}

// audio
func (s *AccessService) audioHandler(userID vo.ID, aggregate agg.Aggregate) (isGranted bool, err error) {
	if !s.videoIsAppropriateHandler(aggregate) {
		return false, s.logger.LogPropagate(
			fmt.Errorf(
				"unable to check access for given aggregate of type '%v' in audio access handler",
				reflect.TypeOf(aggregate).Name(),
			),
		)
	}
	// TODO must be implemented
	return true, nil
}
func (s *AccessService) videoIsAppropriateHandler(aggregate agg.Aggregate) (isAppropriate bool) {
	if _, ok := aggregate.(*agg.Video); ok {
		return true
	}
	return false
}

// resource
func (s *AccessService) resourceHandler(userId vo.ID, aggregate agg.Aggregate) (isGranted bool, err error) {
	if !s.resourceIsAppropriateHandler(aggregate) {
		return false, s.logger.LogPropagate(
			fmt.Errorf(
				"unable to check access for given aggregate of type '%v' in resource access handler",
				reflect.TypeOf(aggregate).Name(),
			),
		)
	}
	// TODO must be implemented
	return true, nil
}
func (s *AccessService) resourceIsAppropriateHandler(v agg.Aggregate) (isAppropriate bool) {
	if _, ok := v.(*agg.Resource); ok {
		return true
	}
	return false
}

// user
func (s *AccessService) userHandler(userID vo.ID, aggregate agg.Aggregate) (isGranted bool, err error) {
	userAgg, ok := aggregate.(*agg.User)
	if !ok {
		return false, s.logger.LogPropagate(
			fmt.Errorf(
				"unable to check access for given aggregate of type '%v' in user access handler",
				reflect.TypeOf(aggregate).Name(),
			),
		)
	}

	if userAgg.ID.Value != userID.Value {
		return false, s.logger.LogPropagate(errors.NewAccessDeniedError("you have not enough rights"))
	}

	return true, nil
}
func (s *AccessService) userIsAppropriateHandler(v agg.Aggregate) (isAppropriate bool) {
	if _, ok := v.(*agg.User); ok {
		return true
	}
	return false
}

func (s *AccessService) setHandlers() *AccessService {
	// video
	s.handlers[VideoType] = s.videoHandler
	s.isAppropriateHandlerFuncs[VideoType] = s.videoIsAppropriateHandler
	// audio
	s.handlers[AudioType] = s.audioHandler
	s.isAppropriateHandlerFuncs[AudioType] = s.audioIsAppropriateHandler
	// resource
	s.handlers[ResourceType] = s.resourceHandler
	s.isAppropriateHandlerFuncs[ResourceType] = s.resourceIsAppropriateHandler
	// user
	s.handlers[UserType] = s.userHandler
	s.isAppropriateHandlerFuncs[UserType] = s.userIsAppropriateHandler
	// fluent setter
	return s
}
