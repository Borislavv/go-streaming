package read

import (
	"fmt"
	"github.com/Borislavv/video-streaming/internal/domain/model"
	"github.com/Borislavv/video-streaming/internal/domain/model/stream"
	"github.com/nareix/joy4/av"
	"github.com/nareix/joy4/av/avutil"
	"github.com/nareix/joy4/format"
	"io"
	"log"
	"os"
	"time"
)

// ChunkSize is 2.5MB
const ChunkSize = 1024 * 1024 * 1

type ReadingService struct {
	errCh chan error
}

func NewReadingService(errCh chan error) *ReadingService {
	return &ReadingService{errCh: errCh}
}

func (r *ReadingService) Read(resource model.Resource) chan *stream.Chunk {
	log.Println("[reader]: reading started")

	chunksCh := make(chan *stream.Chunk, 10)
	go r.handleRead(resource, chunksCh)

	return chunksCh
}

func (r *ReadingService) handleRead(resource model.Resource, chunksCh chan *stream.Chunk) {
	defer log.Println("[reader]: reading stopped")
	defer close(chunksCh)

	file, err := os.Open(resource.GetPath())
	if err != nil {
		r.errCh <- err
		return
	}
	defer func() {
		if err = file.Close(); err != nil {
			r.errCh <- err
			return
		}
	}()

	for {
		chunk := stream.NewChunk(ChunkSize)

		chunk.Len, err = file.Read(chunk.Data)
		if err != nil {
			if err == io.EOF {
				log.Println("[reader]: file was successfully read")
				break
			}
			r.errCh <- err
			return
		}

		if chunk.Len < ChunkSize {
			lastChunk := make([]byte, chunk.Len)
			lastChunk = chunk.Data[:chunk.Len]
			chunk.Data = lastChunk
		}

		if chunk.Len > 0 {
			log.Printf("[reader]: read %d bytes", chunk.Len)
			chunksCh <- chunk
		}

		time.Sleep(time.Second * 1)
	}
}

// todo must be complete as separated module
func DetermineCodecs(resource model.Resource) {
	// Инициализируем библиотеку joy4
	format.RegisterAll()

	// Открываем медиафайл
	f, err := avutil.Open(resource.GetPath())
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	// Получаем первый видео поток из медиафайла
	videoStreams, err := f.Streams()
	if err != nil {
		log.Fatal("Не удалось получить потоки из файла")
	}

	// Отбираем только видео потоки
	var videoStreamsData []av.CodecData
	for _, streamm := range videoStreams {
		if streamm.Type().IsVideo() {
			videoStreamsData = append(videoStreamsData, streamm)
		}
	}

	// Выводим информацию о формате и кодеке видео для всех видео потоков
	for _, videoStream := range videoStreamsData {
		fmt.Printf("Формат видео: %s\n", videoStream.Type())
		fmt.Printf("Кодек видео: %s\n", videoStream.Type().String())
	}
}
