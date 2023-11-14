package entity

import "github.com/Borislavv/video-streaming/internal/domain/vo"

type Video struct {
	ID          vo.ID  `json:"id" bson:",inline"`
	UserID      vo.ID  `json:"userID" bson:"user"`
	Name        string `json:"name" bson:"name"`
	Description string `json:"description" bson:"description,omitempty"`
}

func (r Video) GetID() vo.ID {
	return r.ID
}
