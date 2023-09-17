package reader

type Resource interface {
	GetFilepath() string
	//GetName() string
	//GetFilename() string
	//GetFilesize() int64
	//GetFileMIME() textproto.MIMEHeader
}
