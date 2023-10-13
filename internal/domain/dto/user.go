package dto

import (
	"github.com/Borislavv/video-streaming/internal/domain/vo"
)

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

type UserGetRequestDTO struct {
	ID vo.ID `json:"id"`
}

func (req *UserGetRequestDTO) GetId() vo.ID {
	return req.ID
}
