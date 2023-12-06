package agg

import "github.com/Borislavv/video-streaming/internal/domain/vo"

// Aggregate is interface which must be implemented by all aggregates
// (commonly the target contract will be implemented by embedded identity entity).
type Aggregate interface {
	GetID() vo.ID
}
