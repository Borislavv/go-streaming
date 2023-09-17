package streamer

type Chunk interface {
	GetLen() int
	SetLen(len int)
	GetData() []byte
	SetData(data []byte)
	GetError() error
	SetError(err error)
}
