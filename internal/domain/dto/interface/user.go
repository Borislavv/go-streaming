package dto_interface

import (
	"github.com/Borislavv/video-streaming/internal/domain/vo"
)

type CreateUserRequest interface {
	GetUsername() string
	GetPassword() string
	GetBirthday() string
	GetEmail() string
}

type UpdateUserRequest interface {
	GetID() vo.ID
	GetUsername() string
	GetPassword() string
	GetBirthday() string
}

type GetUserRequest interface {
	GetID() vo.ID
	GetEmail() string
}

type DeleteUserRequest interface {
	GetID() vo.ID
}
