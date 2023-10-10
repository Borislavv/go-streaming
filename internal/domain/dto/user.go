package dto

import "github.com/Borislavv/video-streaming/internal/domain/vo"

type UserCreateRequestDTO struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Birthday string `json:"birthday"`
	Email    string `json:"email"` // unique key
}

type UserUpdateRequestDTO struct {
	ID       vo.ID  `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
	Birthday string `json:"birthday"`
}

type UserGetRequestDTO struct {
	ID vo.ID `json:"id"`
}

func (req *UserGetRequestDTO) GetId() vo.ID {
	return req.ID
}
