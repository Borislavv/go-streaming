package accessorinterface

import (
	dtointerface "github.com/Borislavv/video-streaming/internal/domain/agg/interface"
	"github.com/Borislavv/video-streaming/internal/domain/vo"
)

type Accessor interface {
	// IsGranted is a method which will check the access to target aggregates scope.
	IsGranted(userID vo.ID, aggregates ...dtointerface.Aggregate) error
}
