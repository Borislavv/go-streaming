package entity

import "github.com/Borislavv/video-streaming/internal/domain/vo"

type Video struct {
	ID          vo.ID  `json:"id" bson:",inline"`
	Name        string `json:"name" bson:"name"`
	Description string `json:"description" bson:"description,omitempty"`
}
