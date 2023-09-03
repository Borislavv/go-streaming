package dto

// DefaultChunkSize is one MB
const DefaultChunkSize = 1024 * 1024

type ChunkDto struct {
	Len  int
	Data []byte
	Err  error
}

func NewChunk(size uint64) *ChunkDto {
	if size == 0 {
		size = DefaultChunkSize
	}
	return &ChunkDto{Data: make([]byte, size)}
}

func (c *ChunkDto) GetLen() int {
	return c.Len
}

func (c *ChunkDto) SetLen(len int) {
	c.Len = len
}

func (c *ChunkDto) GetData() []byte {
	return c.Data
}

func (c *ChunkDto) SetData(data []byte) {
	c.Data = data
}

func (c *ChunkDto) GetError() error {
	return c.Err
}

func (c *ChunkDto) SetError(err error) {
	c.Err = err
}
