package reader_interface

import (
	"github.com/Borislavv/video-streaming/internal/infrastructure/service/reader/model"
	"os"
)

type FileReader interface {
	// ReadAll - reads a whole file in a single chunk.
	ReadAll(file *os.File) *model.Chunk
	// ReadByChunks - reads a file by separated chunks and passed it into the channel.
	ReadByChunks(file *os.File, offset int64) chan *model.Chunk
}
