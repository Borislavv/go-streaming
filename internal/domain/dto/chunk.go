package dto

const DefaultChunkSize = 1024 * 1024 // 1MB

type ChunkDTO struct {
	Len  int64
	Data []byte
	Err  error
}

func NewChunk(size uint64) *ChunkDTO {
	if size == 0 {
		size = DefaultChunkSize
	}
	return &ChunkDTO{Data: make([]byte, size)}
}

func (c *ChunkDTO) GetLen() int64 {
	return c.Len
}

func (c *ChunkDTO) SetLen(len int64) {
	c.Len = len
}

func (c *ChunkDTO) GetData() []byte {
	return c.Data
}

func (c *ChunkDTO) SetData(data []byte) {
	c.Data = data
}

func (c *ChunkDTO) GetError() error {
	return c.Err
}

func (c *ChunkDTO) SetError(err error) {
	c.Err = err
}
