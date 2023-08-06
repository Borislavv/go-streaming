package main

import (
	"github.com/Borislavv/video-streaming/cmd/api/resource"
	"github.com/Borislavv/video-streaming/cmd/api/stream"
	"log"
	"sync"
)

func main() {
	wg := &sync.WaitGroup{}
	wg.Add(2)

	// Run streaming service (websocket server)
	go stream.NewApiService().Run(wg)

	// Run static handler service (http server)
	go resource.NewApiService().Run(wg)

	log.Println("[application]: is running")
	wg.Wait()
	log.Println("[application]: was stopped")
}
