package entity

import (
	"github.com/Borislavv/video-streaming/internal/domain/vo"
)

const (
	MIMEContentDispositionKey = "Content-Disposition"
	MIMEContentTypeKey        = "Content-Type"
)

type Resource struct {
	ID       vo.ID  `json:"id" bson:",inline"`
	UserID   vo.ID  `json:"userID" bson:"userID"`     // user identifier
	Name     string `json:"name" bson:"name"`         // original filename
	Filename string `json:"filename" bson:"filename"` // uploaded filename
	Filepath string `json:"filepath" bson:"filepath"` // path to uploaded file
	Filetype string `json:"filetype" bson:"filetype"` // filetype
	Filesize int64  `json:"filesize" bson:"filesize"` // size of uploaded file
}

func (r Resource) GetID() vo.ID {
	return r.ID
}
func (r Resource) GetUserID() vo.ID {
	return r.UserID
}
func (r Resource) GetName() string {
	return r.Name
}
func (r Resource) GetFilename() string {
	return r.Filename
}
func (r Resource) GetFilepath() string {
	return r.Filepath
}
func (r Resource) GetFilesize() int64 {
	return r.Filesize
}
func (r Resource) GetFiletype() string {
	return r.Filetype
}
