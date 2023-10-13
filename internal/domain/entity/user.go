package entity

import (
	"github.com/Borislavv/video-streaming/internal/domain/vo"
	"time"
)

type User struct {
	ID       vo.ID     `json:"id" bson:",inline"`
	Username string    `json:"username" bson:"username"`
	Password string    `json:"password" bson:"password"`
	Email    string    `json:"email" bson:"email"` // unique key
	Birthday time.Time `json:"birthday" bson:"birthday,omitempty"`
}
