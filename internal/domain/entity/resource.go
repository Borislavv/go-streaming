package entity

import "net/textproto"

const (
	MIMEContentDispositionKey = "Content-Disposition"
	MIMEContentTypeKey        = "Content-Type"
)

type Resource struct {
	Name     string               `json:"name" bson:"name"`         // original filename
	Filename string               `json:"filename" bson:"filename"` // uploaded filename
	Path     string               `json:"path" bson:"path"`         // path to uploaded file
	Size     int64                `json:"size" bson:"size"`         // size of file
	MIME     textproto.MIMEHeader `json:"MIME" bson:"MIME"`         // file MIME type
}

func (r *Resource) GetName() string {
	return r.Name
}
func (r *Resource) GetFilename() string {
	return r.Filename
}
func (r *Resource) GetPath() string {
	return r.Path
}
func (r *Resource) GetSize() int64 {
	return r.Size
}
func (r *Resource) GetMIME() textproto.MIMEHeader {
	return r.MIME
}
