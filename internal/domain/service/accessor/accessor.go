package accessor

import (
	"fmt"
	"github.com/Borislavv/video-streaming/internal/domain/agg"
	dtointerface "github.com/Borislavv/video-streaming/internal/domain/agg/interface"
	"github.com/Borislavv/video-streaming/internal/domain/errtype"
	"github.com/Borislavv/video-streaming/internal/domain/logger/interface"
	"github.com/Borislavv/video-streaming/internal/domain/service/di/interface"
	"github.com/Borislavv/video-streaming/internal/domain/vo"
	"reflect"
)

type AggregateAccessType int
type AggregateAccessHandler func(userID vo.ID, agg dtointerface.Aggregate) (err error)
type AggregateAccessIsAppropriateHandler func(agg dtointerface.Aggregate) (isSupported bool)

const (
	VideoType AggregateAccessType = iota
	AudioType
	ResourceType
	UserType
)

type AccessService struct {
	logger                    loggerinterface.Logger
	handlers                  map[AggregateAccessType]AggregateAccessHandler
	isAppropriateHandlerFuncs map[AggregateAccessType]AggregateAccessIsAppropriateHandler
}

func NewAccessService(serviceContainer diinterface.ContainerManager) (*AccessService, error) {
	loggerService, err := serviceContainer.GetLoggerService()
	if err != nil {
		return nil, err
	}

	return (&AccessService{
		logger:                    loggerService,
		handlers:                  map[AggregateAccessType]AggregateAccessHandler{},
		isAppropriateHandlerFuncs: map[AggregateAccessType]AggregateAccessIsAppropriateHandler{},
	}).setHandlers(), nil
}

// IsGranted is a method which will check the access to target scope of aggregates.
func (s *AccessService) IsGranted(userID vo.ID, aggregates ...dtointerface.Aggregate) error {
	for _, aggregate := range aggregates {
		for aggregateAccessType, isAppropriateHandler := range s.isAppropriateHandlerFuncs {
			if isAppropriateHandler(aggregate) {
				appropriateHandler, found := s.handlers[aggregateAccessType]
				if !found {
					return s.logger.LogPropagate(
						fmt.Errorf(
							"appropriate handler was not found by type '%d' into the handlers map",
							aggregateAccessType,
						),
					)
				}

				if err := appropriateHandler(userID, aggregate); err != nil {
					return s.logger.LogPropagate(err)
				}
			}
		}
	}

	// access is granted, no errors were occurred
	return nil
}

// video
func (s *AccessService) videoHandler(userID vo.ID, aggregate dtointerface.Aggregate) error {
	videoAgg, ok := aggregate.(*agg.Video)
	if !ok {
		return s.logger.LogPropagate(
			fmt.Errorf(
				"unable to check access for given aggregate of type '%v' in video access handler",
				reflect.TypeOf(aggregate).Name(),
			),
		)
	}

	if userID.Value == videoAgg.UserID.Value && userID.Value == videoAgg.Resource.UserID.Value {
		// user is owner of video
		return nil
	}

	// video was not matched, access is denied
	return s.logger.LogPropagate(
		errtype.NewAccessDeniedError(
			"you have not enough rights, one of entities is video and it's not belong to you",
		),
	)
}
func (s *AccessService) audioIsAppropriateHandler(aggregate dtointerface.Aggregate) (isAppropriate bool) {
	if _, ok := aggregate.(*agg.Audio); ok {
		return true
	}
	return false
}

// audio
func (s *AccessService) audioHandler(userID vo.ID, aggregate dtointerface.Aggregate) error {
	audioAgg, ok := aggregate.(*agg.Audio)
	if !ok {
		return s.logger.LogPropagate(
			fmt.Errorf(
				"unable to check access for given aggregate of type '%v' in audio access handler",
				reflect.TypeOf(aggregate).Name(),
			),
		)
	}

	if userID.Value == audioAgg.UserID.Value {
		// user is owner of audio
		return nil
	}

	// audio was not matched, access is denied
	return s.logger.LogPropagate(
		errtype.NewAccessDeniedError(
			"you have not enough rights, one of entities is audio and it's not belong to you",
		),
	)
}
func (s *AccessService) videoIsAppropriateHandler(aggregate dtointerface.Aggregate) (isAppropriate bool) {
	if _, ok := aggregate.(*agg.Video); ok {
		return true
	}
	return false
}

// resource
func (s *AccessService) resourceHandler(userID vo.ID, aggregate dtointerface.Aggregate) error {
	resourceAgg, ok := aggregate.(*agg.Resource)
	if !ok {
		return s.logger.LogPropagate(
			fmt.Errorf(
				"unable to check access for given aggregate of type '%v' in resource access handler",
				reflect.TypeOf(aggregate).Name(),
			),
		)
	}

	if userID.Value == resourceAgg.UserID.Value {
		// user is owner of resource
		return nil
	}

	// resource was not matched, access is denied
	return s.logger.LogPropagate(
		errtype.NewAccessDeniedError(
			"you have not enough rights, one of entities is resource and it's not belong to you",
		),
	)
}
func (s *AccessService) resourceIsAppropriateHandler(v dtointerface.Aggregate) (isAppropriate bool) {
	if _, ok := v.(*agg.Resource); ok {
		return true
	}
	return false
}

// user
func (s *AccessService) userHandler(userID vo.ID, aggregate dtointerface.Aggregate) error {
	userAgg, ok := aggregate.(*agg.User)
	if !ok {
		return s.logger.LogPropagate(
			fmt.Errorf(
				"unable to check access for given aggregate of type '%v' in user access handler",
				reflect.TypeOf(aggregate).Name(),
			),
		)
	}

	if userID.Value == userAgg.ID.Value {
		// users are equal
		return nil
	}

	// user was not matched, access is denied
	return s.logger.LogPropagate(
		errtype.NewAccessDeniedError(
			"you have not enough rights, one of entities is another user",
		),
	)
}
func (s *AccessService) userIsAppropriateHandler(v dtointerface.Aggregate) (isAppropriate bool) {
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
