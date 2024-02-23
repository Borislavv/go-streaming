package dtointerface

type Chunk interface {
	GetLen() int
	GetData() []byte
	SetData(data []byte)
	GetError() error
	SetError(err error)
}
