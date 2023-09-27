package reader

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"github.com/Borislavv/video-streaming/internal/domain/dto"
	"github.com/Borislavv/video-streaming/internal/domain/logger"
	"io"
	"math"
	"os"
	"sync"
)

const (
	ChunkSize      = 1024 * 1024 * 1 // 1MB
	ChunksBuffer   = 4
	ReadingThreads = 4
)

type ResourceReader struct {
	logger logger.Logger
	cache  map[string]map[int]dto.Chunk
}

func NewReaderService(logger logger.Logger) *ResourceReader {
	return &ResourceReader{
		logger: logger,
		cache:  map[string]map[int]dto.Chunk{},
	}
}

// Read will read a resource and send file as butches of bytes
func (r *ResourceReader) Read(resource dto.Resource) chan dto.Chunk {
	r.logger.Info(fmt.Sprintf("recourse '%v' reading started", resource.GetName()))

	chunksCh := make(chan dto.Chunk, ChunksBuffer)
	go r.handleRead(resource, chunksCh)

	return chunksCh
}

func (r *ResourceReader) handleRead(resource dto.Resource, chunksCh chan dto.Chunk) {
	defer func() {
		close(chunksCh)
		r.logger.Info(fmt.Sprintf("recourse '%v' reading finished", resource.GetName()))
	}()

	chunkedFile, err := r.cached(resource)
	if err != nil {
		r.logger.Emergency(err)
		return
	}

	for i := 0; i < len(chunkedFile); i++ {
		r.sendChunk(chunkedFile[i], chunksCh)
	}

	// todo must be separated to own native: single thread strategy
	//file, err := os.Open(resource.GetFilepath())
	//if err != nil {
	//	r.logger.Error(err)
	//	return
	//}
	//defer func() { _ = file.Close() }()
	//
	//for {
	//	chunk := dto.NewChunk(ChunkSize, 0)
	//
	//	chunk.Len, err = file.Read(chunk.Data)
	//	if err != nil {
	//		if err == io.EOF {
	//			break
	//		}
	//		r.logger.Error(err)
	//		return
	//	}
	//
	//	r.sendChunk(chunk, chunksCh)
	//}
}

func (r *ResourceReader) sendChunk(chunk dto.Chunk, chunksCh chan dto.Chunk) {
	if chunk.GetLen() < ChunkSize {
		lastChunk := make([]byte, chunk.GetLen())
		lastChunk = chunk.GetData()[:chunk.GetLen()]
		chunk.SetData(lastChunk)
	}

	if chunk.GetLen() > 0 {
		chunksCh <- chunk
	}
}

func (r *ResourceReader) read(resource dto.Resource) map[int]dto.Chunk {
	file, err := os.Open(resource.GetFilepath())
	if err != nil {
		r.logger.Error(err)
		return nil
	}
	defer func() { _ = file.Close() }()

	info, err := file.Stat()
	if err != nil {
		r.logger.Error(err)
		return nil
	}
	// computing number of chunks for read full file
	chunksNum := int(math.Ceil(float64(info.Size()) / ChunkSize))
	// computing number of active reading threads
	threads := ReadingThreads
	if chunksNum < ReadingThreads {
		threads = chunksNum
	}

	type ChunkOffset struct {
		Number int
		Offset int64
	}

	// initialize cache blocks
	cache := map[int]dto.Chunk{}

	wgThreads := &sync.WaitGroup{}
	offsetCh := make(chan ChunkOffset, threads)

	wgThreads.Add(1)
	go func() {
		defer func() {
			close(offsetCh)
			wgThreads.Done()
		}()
		for chk := 0; chk < chunksNum; chk++ {
			offsetCh <- ChunkOffset{
				Number: chk,
				Offset: int64(ChunkSize * chk),
			}
		}
	}()

	chunksCh := make(chan dto.Chunk, threads)
	wgThreads.Add(threads)
	for thr := 0; thr < threads; thr++ {
		go func() {
			defer wgThreads.Done()

			for offset := range offsetCh {
				// building a new chunk
				chunk := dto.NewChunk(ChunkSize, offset.Number)
				// reading with offset
				chunk.Len, err = file.ReadAt(chunk.Data, offset.Offset)
				if err != nil {
					if err == io.EOF {
						// sent the last chunk, if it is not empty
						if chunk.GetLen() > 0 {
							chunksCh <- chunk
						}
						return
					}
					r.logger.Error(err)
					return
				}
				// thread safety due to use of work provider
				chunksCh <- chunk
			}
		}()
	}

	wgConsume := &sync.WaitGroup{}
	wgConsume.Add(1)
	go func() {
		defer wgConsume.Done()

		for chunk := range chunksCh {
			cache[chunk.GetNum()] = chunk
		}
	}()

	wgThreads.Wait()
	close(chunksCh)
	wgConsume.Wait()

	// store to cache
	return cache
}

// TODO must be implemented: cache eviction!
func (r *ResourceReader) cached(resource dto.Resource) (map[int]dto.Chunk, error) {
	hash := md5.New()
	if _, err := hash.Write([]byte(resource.GetFilepath())); err != nil {
		r.logger.Emergency(err)
		return nil, err
	}
	key := hex.EncodeToString(hash.Sum(nil))

	// if data is found
	if data, found := r.cache[key]; found {
		r.logger.Info(fmt.Sprintf("resource '%v' was fetched from cache", resource.GetName()))
		return data, nil
	}

	// if data is not found
	data := r.read(resource)
	r.cache[key] = data

	r.logger.Info(fmt.Sprintf("resource '%v' was fetched from file storage", resource.GetName()))

	return data, nil
}
