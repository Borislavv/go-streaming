package query

import "github.com/Borislavv/video-streaming/internal/domain/vo"

type HasBlockedToken struct {
	token  string
	userID vo.ID
}

func NewHasBlockedToken(token string, userID vo.ID) *HasBlockedToken {
	return &HasBlockedToken{
		token:  token,
		userID: userID,
	}
}
func (q *HasBlockedToken) GetToken() string {
	return q.token
}
func (q *HasBlockedToken) GetUserID() vo.ID {
	return q.userID
}
