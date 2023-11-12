package query

import (
	"github.com/Borislavv/video-streaming/internal/domain/vo"
	"time"
)

type FindOneVideoByID struct {
	id     vo.ID
	userId vo.ID
}

func NewFindOneVideoByID(id vo.ID, userId vo.ID) *FindOneVideoByID {
	return &FindOneVideoByID{
		id:     id,
		userId: userId,
	}
}
func (q FindOneVideoByID) GetID() vo.ID {
	return q.id
}
func (q FindOneVideoByID) GetUserID() vo.ID {
	return q.userId
}

type FindOneVideoByName struct {
	name   string
	userID vo.ID
}

func NewFindOneVideoByName(name string, userID vo.ID) *FindOneVideoByName {
	return &FindOneVideoByName{
		name:   name,
		userID: userID,
	}
}
func (q *FindOneVideoByName) GetName() string {
	return q.name
}
func (q *FindOneVideoByName) GetUserID() vo.ID {
	return q.userID
}

type FindOneVideoByResourceID struct {
	resourceID vo.ID
	userID     vo.ID
}

func NewFindOneVideoByResourceID(resourceID vo.ID, userID vo.ID) *FindOneVideoByResourceID {
	return &FindOneVideoByResourceID{
		resourceID: resourceID,
		userID:     userID,
	}
}
func (q *FindOneVideoByResourceID) GetResourceID() vo.ID {
	return q.resourceID
}
func (q *FindOneVideoByResourceID) GetUserID() vo.ID {
	return q.userID
}

type FindVideoList struct {
	name      string    // part of name
	userID    vo.ID     // user identifier
	createdAt time.Time // concrete search date point
	from      time.Time // search date limit from
	to        time.Time // search date limit to
	Paginated
}

func NewFindVideoList(
	name string,
	userID vo.ID,
	createdAt time.Time,
	from time.Time,
	to time.Time,
	page int,
	limit int,
) *FindVideoList {
	return &FindVideoList{
		name:      name,
		userID:    userID,
		createdAt: createdAt,
		from:      from,
		to:        to,
		Paginated: Paginated{
			page:  page,
			limit: limit,
		},
	}
}
func (q *FindVideoList) GetName() string {
	return q.name
}
func (q *FindVideoList) GetUserID() vo.ID {
	return q.userID
}
func (q *FindVideoList) GetCreatedAt() time.Time {
	return q.createdAt
}
func (q *FindVideoList) GetFrom() time.Time {
	return q.from
}
func (q *FindVideoList) GetTo() time.Time {
	return q.to
}
