package main

import (
	"github.com/Borislavv/video-streaming/internal/app/resource"
	"github.com/Borislavv/video-streaming/internal/app/stream"
	"sync"
)

func main() {
	wg := &sync.WaitGroup{}
	wg.Add(2)

	// Run streaming app (websocket server)
	go stream.NewStreamingApp().Run(wg)

	// RestApi, Static files serving, Native rendering (http server)
	go resource.NewResourcesApp().Run(wg)

	wg.Wait()
}
