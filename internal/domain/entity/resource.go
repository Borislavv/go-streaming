package entity

import (
	"github.com/Borislavv/video-streaming/internal/domain/vo"
	"net/textproto"
)

const (
	MIMEContentDispositionKey = "Content-Disposition"
	MIMEContentTypeKey        = "Content-Type"
)

type Resource struct {
	ID       vo.ID                `json:"id" bson:",inline"`
	Name     string               `json:"name" bson:"name"`         // original filename
	Filename string               `json:"filename" bson:"filename"` // uploaded filename
	Filepath string               `json:"path" bson:"path"`         // path to uploaded file
	Filesize int64                `json:"size" bson:"size"`         // size of uploaded file
	FileMIME textproto.MIMEHeader `json:"MIME" bson:"MIME"`         // file MIME type
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
func (r Resource) GetFileMIME() textproto.MIMEHeader {
	return r.FileMIME
}
func (r Resource) GetContentType() string {
	return r.GetFileMIME().Get(MIMEContentTypeKey)
}
func (r Resource) GetContentDisposition() string {
	return r.GetFileMIME().Get(MIMEContentDispositionKey)
}
