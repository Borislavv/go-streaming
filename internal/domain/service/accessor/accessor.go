package accessor

import (
	"context"
	"fmt"
	"github.com/Borislavv/video-streaming/internal/domain/agg"
	"github.com/Borislavv/video-streaming/internal/domain/errors"
	"github.com/Borislavv/video-streaming/internal/domain/logger"
	"github.com/Borislavv/video-streaming/internal/domain/vo"
	"reflect"
)

type AggregateAccessType int
type AggregateAccessHandler func(userID vo.ID, agg agg.Aggregate) (err error)
type AggregateAccessIsAppropriateHandler func(agg agg.Aggregate) (isSupported bool)

const (
	VideoType AggregateAccessType = iota
	AudioType
	ResourceType
	UserType
)

type AccessService struct {
	logger                    logger.Logger
	handlers                  map[AggregateAccessType]AggregateAccessHandler
	isAppropriateHandlerFuncs map[AggregateAccessType]AggregateAccessIsAppropriateHandler
}

func NewAccessService(ctx context.Context, logger logger.Logger) *AccessService {
	return (&AccessService{
		logger:                    logger,
		handlers:                  map[AggregateAccessType]AggregateAccessHandler{},
		isAppropriateHandlerFuncs: map[AggregateAccessType]AggregateAccessIsAppropriateHandler{},
	}).setHandlers()
}

// IsGranted is a method which will check the access to target scope of aggregates.
func (s *AccessService) IsGranted(userID vo.ID, aggregates ...agg.Aggregate) error {
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
func (s *AccessService) videoHandler(userID vo.ID, aggregate agg.Aggregate) error {
	videoAgg, ok := aggregate.(*agg.Video)
	if !ok {
		return s.logger.LogPropagate(
			fmt.Errorf(
				"unable to check access for given aggregate of type '%v' in video access handler",
				reflect.TypeOf(aggregate).Name(),
			),
		)
	}

	if userID.Value == videoAgg.UserID.Value {
		if err := s.resourceHandler(userID, videoAgg.Resource); err != nil {
			return s.logger.LogPropagate(err)
		}
		// user is owner of video
		return nil
	}

	// video was not matched, access is denied
	return s.logger.LogPropagate(
		errors.NewAccessDeniedError(
			"you have not enough rights, one of entities is video and it's not belong to you",
		),
	)
}
func (s *AccessService) audioIsAppropriateHandler(aggregate agg.Aggregate) (isAppropriate bool) {
	if _, ok := aggregate.(*agg.Audio); ok {
		return true
	}
	return false
}

// audio
func (s *AccessService) audioHandler(userID vo.ID, aggregate agg.Aggregate) error {
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
		errors.NewAccessDeniedError(
			"you have not enough rights, one of entities is audio and it's not belong to you",
		),
	)
}
func (s *AccessService) videoIsAppropriateHandler(aggregate agg.Aggregate) (isAppropriate bool) {
	if _, ok := aggregate.(*agg.Video); ok {
		return true
	}
	return false
}

// resource
func (s *AccessService) resourceHandler(userID vo.ID, aggregate agg.Aggregate) error {
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
		errors.NewAccessDeniedError(
			"you have not enough rights, one of entities is resource and it's not belong to you",
		),
	)
}
func (s *AccessService) resourceIsAppropriateHandler(v agg.Aggregate) (isAppropriate bool) {
	if _, ok := v.(*agg.Resource); ok {
		return true
	}
	return false
}

// user
func (s *AccessService) userHandler(userID vo.ID, aggregate agg.Aggregate) error {
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
		errors.NewAccessDeniedError(
			"you have not enough rights, one of entities is another user",
		),
	)
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
