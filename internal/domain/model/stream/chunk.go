package stream

// DefaultChunkSize is one MB
const DefaultChunkSize = 1024 * 1024

type Chunk struct {
	Len  int
	Data []byte
	Err  error
}

func NewChunk(size uint64) *Chunk {
	if size == 0 {
		size = DefaultChunkSize
	}
	return &Chunk{Data: make([]byte, size)}
}
