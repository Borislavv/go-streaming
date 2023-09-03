package main

import (
	"github.com/Borislavv/video-streaming/internal/app/resource"
	"github.com/Borislavv/video-streaming/internal/app/stream"
	"sync"
)

func main() {
	wg := &sync.WaitGroup{}
	wg.Add(2)

	// Run streaming service (websocket server)
	go stream.NewStreamingApp().Run(wg)

	// Run static handler service (http server)
	go resource.NewResourcesApp().Run(wg)

	wg.Wait()
}
