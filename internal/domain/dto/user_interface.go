package dto

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
	GetId() vo.ID
}

type DeleteUserRequest GetUserRequest
