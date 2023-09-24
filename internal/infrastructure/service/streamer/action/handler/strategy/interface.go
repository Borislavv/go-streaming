package strategy

import (
	"github.com/Borislavv/video-streaming/internal/infrastructure/service/streamer/action/model"
)

type ActionStrategy interface {
	IsAppropriate(action model.Action) bool
	Do(action model.Action) error
}
