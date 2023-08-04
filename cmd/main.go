package main

import (
	"github.com/Borislavv/video-streaming/cmd/api/resource"
	"github.com/Borislavv/video-streaming/cmd/api/stream"
	"sync"
)

func main() {
	wg := &sync.WaitGroup{}
	wg.Add(2)

	// Run streaming service (websocket server)
	go stream.NewApiService().Run(wg)

	// Run resource handler service (http server)
	go resource.NewApiService().Run(wg)

	wg.Wait()
}
