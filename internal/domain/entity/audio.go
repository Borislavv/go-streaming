package entity

import "github.com/Borislavv/video-streaming/internal/domain/vo"

type Audio struct {
	ID          vo.ID  `json:"id,omitempty" bson:"_id,omitempty,inline"`
	UserID      vo.ID  `json:"userID" bson:"userID"`
	Name        string `bson:"name"`
	Description string `bson:"description,omitempty"`
}

func (r Audio) GetID() vo.ID {
	return r.ID
}
