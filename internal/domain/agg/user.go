package agg

import (
	"github.com/Borislavv/video-streaming/internal/domain/entity"
	"github.com/Borislavv/video-streaming/internal/domain/vo"
)

type User struct {
	entity.User `bson:",inline"`

	VideoIDs    []vo.ID `json:"videos" bson:"videos"`
	AudioIDs    []vo.ID `json:"audios" bson:"audios"`
	ResourceIDs []vo.ID `json:"resources" bson:"resources"`

	Timestamp vo.Timestamp `json:"timestamp" bson:",inline"`
}
