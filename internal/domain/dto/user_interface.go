package dto

import (
	"github.com/Borislavv/video-streaming/internal/domain/vo"
	"time"
)

type CreateUserRequest interface {
	GetUsername() string
	GetPassword() string
	GetBirthday() time.Time
	GetEmail() string
}

type UpdateUserRequest interface {
	GetUsername() string
	GetPassword() string
	GetBirthday() time.Time
	GetEmail() string
}

type GetUserRequest interface {
	GetId() vo.ID
}

type DeleteUserRequest GetUserRequest
