package main

import (
	"github.com/Borislavv/video-streaming/cmd/api/resource"
	"github.com/Borislavv/video-streaming/cmd/api/stream"
)

func main() {
	// Run streaming service (websocket server)
	go stream.NewApiService().Run()

	// Run resource handler service (http server)
	resource.NewApiService().Run()
}
