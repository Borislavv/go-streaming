package query

import "github.com/Borislavv/video-streaming/internal/domain/vo"

type FindOneUserByID struct {
	id vo.ID
}

func NewFindOneUserByID(id vo.ID) *FindOneUserByID {
	return &FindOneUserByID{id: id}
}
func (q *FindOneUserByID) GetID() vo.ID {
	return q.id
}

type FindOneUserByEmail struct {
	email string
}

func NewFindOneUserByEmail(email string) *FindOneUserByEmail {
	return &FindOneUserByEmail{email: email}
}
func (q *FindOneUserByEmail) GetEmail() string {
	return q.email
}
