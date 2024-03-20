package queryinterface

import "github.com/Borislavv/video-streaming/internal/domain/vo"

type FindOneUserByID interface {
	GetID() vo.ID
}

type FindOneUserByEmail interface {
	GetEmail() string
}
