package dto

import (
	"github.com/Borislavv/video-streaming/internal/domain/vo"
	"time"
)

// UserCreateRequestDTO - used when u want to create the user
type UserCreateRequestDTO struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Birthday string `json:"birthday"`
	Email    string `json:"email"` // unique key
}

func (u *UserCreateRequestDTO) GetUsername() string {
	return u.Username
}
func (u *UserCreateRequestDTO) GetPassword() string {
	return u.Password
}
func (u *UserCreateRequestDTO) GetBirthday() string {
	return u.Birthday
}
func (u *UserCreateRequestDTO) GetEmail() string {
	return u.Email
}

// UserUpdateRequestDTO - used when u want to update the user
type UserUpdateRequestDTO struct {
	ID       vo.ID  `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
	Birthday string `json:"birthday"`
}

func (u *UserUpdateRequestDTO) GetID() vo.ID {
	return u.ID
}
func (u *UserUpdateRequestDTO) GetUsername() string {
	return u.Username
}
func (u *UserUpdateRequestDTO) GetPassword() string {
	return u.Password
}
func (u *UserUpdateRequestDTO) GetBirthday() string {
	return u.Birthday
}

// UserGetRequestDTO - used when u want to get the user
type UserGetRequestDTO struct {
	ID    vo.ID  `json:"id"`
	Email string `json:"email"`
}

func NewUserGetRequestDTO(id vo.ID, email string) *UserGetRequestDTO {
	return &UserGetRequestDTO{ID: id, Email: email}
}

func (req *UserGetRequestDTO) GetID() vo.ID {
	return req.ID
}

func (req *UserGetRequestDTO) GetEmail() string {
	return req.Email
}

// UserDeleteRequestDTO - used when u want to remove the user
type UserDeleteRequestDTO struct {
	ID vo.ID `json:"id"`
}

func (req *UserDeleteRequestDTO) GetID() vo.ID {
	return req.ID
}

type UserResponseDTO struct {
	ID       vo.ID     `json:"ID" bson:",inline"`
	Username string    `json:"username" bson:"username"`
	Email    string    `json:"email" bson:"email"` // unique key
	Birthday time.Time `json:"birthday" bson:"birthday,omitempty"`
}
