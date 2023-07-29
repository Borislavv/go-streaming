package model

type Chunk struct {
	Len  int
	Data []byte
	Err  error
}
