package query_interface

import "github.com/Borislavv/video-streaming/internal/domain/vo"

type FindOneResourceByID interface {
	GetID() vo.ID
	GetUserID() vo.ID
}
