package main

import (
	"github.com/Borislavv/video-streaming/internal/app/resource"
	"github.com/Borislavv/video-streaming/internal/app/stream"
	"github.com/Borislavv/video-streaming/internal/domain/service/di"
	"sync"
)

func main() {
	wg := &sync.WaitGroup{}
	wg.Add(2)

	// Run streaming app (websocket server)
	go stream.NewStreamingApp(di.NewServiceContainerManager()).Run(wg)

	// RestApi, Static files serving, Native rendering (http server)
	go resource.NewResourcesApp(di.NewServiceContainerManager()).Run(wg)

	wg.Wait()
}
