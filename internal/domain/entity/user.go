package entity

import "github.com/Borislavv/video-streaming/internal/domain/vo"

type User struct {
	ID       vo.ID  `json:"ID" bson:",inline"`
	Username string `json:"username" bson:"username"`
	Password string `json:"password" bson:"password"`
	Birthday string `json:"birthday" bson:"birthday,omitempty"`
	Email    string `json:"email" bson:"email"` // unique key
}
