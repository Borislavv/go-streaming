package model

type Chunk struct {
	Len  int64
	Data []byte
	Err  error
}

func NewChunk(size int64) *Chunk {
	return &Chunk{Data: make([]byte, size)}
}

func (c *Chunk) GetLen() int64 {
	return c.Len
}

func (c *Chunk) SetLen(len int64) {
	c.Len = len
}

func (c *Chunk) GetData() []byte {
	return c.Data
}

func (c *Chunk) SetData(data []byte) {
	c.Data = data
}

func (c *Chunk) GetError() error {
	return c.Err
}

func (c *Chunk) SetError(err error) {
	c.Err = err
}
