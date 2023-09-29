package model

import "io"

type Chunk struct {
	Data []byte
	Err  error
}

func NewChunk(length int64, capacity int64) *Chunk {
	return &Chunk{Data: make([]byte, length, capacity)}
}

func (c *Chunk) Read(p []byte) (n int, err error) {
	if c.Data == nil || len(c.Data) == 0 {
		return 0, io.EOF
	}
	n = copy(p, c.Data)
	c.Data = c.Data[n:]
	return n, nil
}
func (c *Chunk) GetLen() int {
	return len(c.Data)
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
