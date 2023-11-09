package accessor

import (
	"github.com/Borislavv/video-streaming/internal/domain/agg"
	"github.com/Borislavv/video-streaming/internal/domain/vo"
)

type Accessor interface {
	// IsGranted is a method which will check the access to target aggregates scope.
	IsGranted(userID vo.ID, aggregates ...agg.Aggregate) error
}
