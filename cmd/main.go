package main

import (
	"github.com/Borislavv/video-streaming/cmd/api/resource"
)

func main() {
	// Run streaming service (websocket server)
	//stream.NewApiService().Run()

	// Run resource handler service (http server)
	resource.NewApiService().Run()
}
