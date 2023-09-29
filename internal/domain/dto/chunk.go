package dto

type ChunkDTO struct {
	Data []byte
	Err  error
}

func NewChunk(length int64, capacity uint64) *ChunkDTO {
	return &ChunkDTO{Data: make([]byte, length, capacity)}
}

func (c *ChunkDTO) GetLen() int {
	return len(c.Data)
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
