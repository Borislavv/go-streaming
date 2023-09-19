package entity

import "github.com/Borislavv/video-streaming/internal/domain/vo"

type Audio struct {
	ID          vo.ID  `json:"id,omitempty" bson:"_id,omitempty,inline"`
	Name        string `bson:"name"`
	Description string `bson:"description,omitempty"`
}
