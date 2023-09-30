package helper

import (
	"github.com/Borislavv/video-streaming/internal/domain/errors"
	"time"
)

var dateLayouts = []string{
	"2006-01-02T15:04:05Z07:00",
	"2006-01-02T15:04:05",
	"2006-01-02",
}

// ParseTime is a user helper function which parsing three types of datetime strings:
// 	1. 2006-01-02T15:04:05Z07:00
//	2. 2006-01-02T15:04:05
// 	3. 2006-01-02
func ParseTime(date string) (time.Time, error) {
	for _, layout := range dateLayouts {
		parsed, err := time.Parse(layout, date)
		if err != nil {
			continue
		}
		return parsed, nil
	}
	return time.Time{}, errors.NewTimeParsingValidationError(date)
}
