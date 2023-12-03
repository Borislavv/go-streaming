package entity

import (
	"github.com/Borislavv/video-streaming/internal/domain/vo"
	"time"
)

type User struct {
	ID       vo.ID     `bson:",inline"`
	Username string    `bson:"username"`
	Password string    `bson:"password"` // hash
	Email    string    `bson:"email"`    // unique key
	Birthday time.Time `bson:"birthday,omitempty"`
}

func (r *User) GetID() vo.ID {
	return r.ID
}
func (r *User) GetPassword() string {
	return r.Password
}
func (r *User) SetPassword(password string) {
	r.Password = password
}
