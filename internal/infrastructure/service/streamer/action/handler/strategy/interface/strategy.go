package strategyinterface

import (
	"github.com/Borislavv/video-streaming/internal/infrastructure/service/streamer/action/model"
)

type ActionStrategy interface {
	// IsAppropriate - method will tell the service architect that the strategy is acceptable.
	IsAppropriate(action model.Action) bool
	// Do is a method which contains the useful work of target strategy.
	Do(action model.Action) error
}
