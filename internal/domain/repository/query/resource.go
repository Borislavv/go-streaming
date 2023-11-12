package query

import "github.com/Borislavv/video-streaming/internal/domain/vo"

type FindOneResourceByID struct {
	id     vo.ID
	userID vo.ID
}

func NewFindOneResourceByID(id vo.ID, userID vo.ID) *FindOneResourceByID {
	return &FindOneResourceByID{
		id:     vo.ID{},
		userID: vo.ID{},
	}
}
func (q *FindOneResourceByID) GetID() vo.ID {
	return q.id
}
func (q *FindOneResourceByID) GetUserID() vo.ID {
	return q.userID
}
