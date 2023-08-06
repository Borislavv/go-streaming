package file

type Chunk struct {
	Len   int
	Bytes []byte
}

func NewChunk(size uint64) *Chunk {
	return &Chunk{Bytes: make([]byte, size)}
}
